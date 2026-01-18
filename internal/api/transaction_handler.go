package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	//"transaction-manager/internal/config"

	"transaction-manager/internal/config"
	"transaction-manager/internal/domain"
	"transaction-manager/internal/repository"
	"transaction-manager/internal/service"
	"transaction-manager/pkg/logger"
)

type TransactionHandler struct {
	service service.TransactionService
	event   service.AccountEventPublisher
}

func NewTransactionHandler(s service.TransactionService, e service.AccountEventPublisher) *TransactionHandler {
	return &TransactionHandler{service: s, event: e}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	reqID := c.GetString("request_id")

	log := logger.Log.With(
		zap.String("service", "transaction-service"),
		zap.String("flow", "CreateTransaction"),
		zap.String("request_id", reqID),
	)
	log.Info("REQUEST_RECEIVED")
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Info("REQUEST_VALIDATED")
	tx := &domain.Transaction{
		UserRefId:        req.UserRefId,
		SourceRefId:      req.SourceRefId,
		DestinationRefId: req.DestinationRefId,
		DcFlag:           req.DcFlag,
		Status:           domain.StatusInitiated,
		PaymentType:      req.PaymentType,
		PaymentMode:      req.PaymentMode,
		SagaStatus:       string(domain.SagaInit),
		Amount:           req.Amount,
		Currency:         req.Currency,
		NetworkTxnId:     req.NetworkTxnId,
		GatewayTxnId:     req.GatewayTxnId,
	}

	var err error
	log.Info("BUSINESS_LOGIC_STARTED")

	if req.PaymentMode == "NEFT" {
		err = h.service.CreateNEFTTransaction(tx)
	} else {
		err = h.service.CreateImmediateTransaction(tx)
	}

	if err != nil {
		log.Error("BUSINESS_LOGIC_FAILED",
			zap.String("reason", "rule_violation"),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// emit event to account service to request balance block

	go h.event.PublishToAccountService(tx, string(domain.StatusBlockRequested), config.KafkaAccountBalanceBlockCmd, context.Background())
	h.service.RecordSagaStep(tx, string(domain.SagaBalanceBlocked), domain.SagaStatusInProgress)
	// DEBIT_REQUESTED	REQUESTED
	h.service.UpdateTransactionWithSaga(tx, domain.StatusBlockRequested, string(domain.SagaBalanceBlocked)+"_"+domain.SagaStatusInProgress)

	resp := CreateTransactionResponse{
		TransactionID: tx.ID,
		Status:        string(domain.StatusPending),
		Message:       "Transaction is being processed",
	}
	log.Info("RESPONSE_SENT",
		zap.Int("status", 201),
	)
	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {

	id := c.Param("id")

	tx, ledger, err := h.service.GetTransaction(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var ledgerResp []LedgerEntryResponse
	for _, l := range ledger {
		ledgerResp = append(ledgerResp, LedgerEntryResponse{
			EntryType: string(l.EntryType),
			Amount:    l.Amount,
			Source:    l.Source,
		})
	}

	resp := GetTransactionResponse{
		TransactionID: tx.ID,
		Status:        string(tx.Status),
		PaymentMode:   tx.PaymentMode,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Ledger:        ledgerResp,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) HandleAccBalBlocked(msg []byte, eventRepo repository.PgxEventRepo) {
	logger.Log.Info("Handling Account Balance Blocked event in TM", zap.ByteString("message", msg))

	// Idempotent processing
	var event domain.ConsumerEventForAccountService
	err := json.Unmarshal(msg, &event)
	if err != nil {
		logger.Log.Error("EVENT_UNMARSHAL_FAILED", zap.Error(err))
		return
	}

	processed, err := eventRepo.IsProcessed(event.EventID)
	if err != nil {
		logger.Log.Error("IDEMPOTENCY_CHECK_FAILED", zap.Error(err))
		return
	}
	if processed {
		logger.Log.Info("DUPLICATE_EVENT_IGNORED - Update Tx and SagaState", zap.String("event_id", event.EventID))
		return
	}

	logger.Log.Info("PROCESSING_ACCOUNT_BALANCE_BLOCKED_EVENT",
		zap.String("event_id", event.EventID),
		zap.String("event_type", event.EventType),
	)
	Tx, _, err := h.service.GetTransaction(event.TransactionID)
	if err != nil {
		logger.Log.Error("TRANSACTION_NOT_FOUND_FOR_EVENT",
			zap.String("event_id", event.EventID),
			zap.String("tx_id", event.TransactionID),
		)
		return
	}
	h.service.RecordSagaStep(Tx, string(domain.SagaBalanceBlocked), domain.SagaStatusCompleted)
	h.service.UpdateTransactionWithSaga(Tx, domain.StatusNetworkRequested, domain.SagaStatusInProgress)

	go h.event.PublishToAccountService(Tx, string(domain.StatusNetworkRequested), config.KafkaPaymentIMPSDebitCmd, context.Background())
	h.service.RecordSagaStep(Tx, string(domain.StatusNetworkRequested), domain.SagaStatusInProgress)
}

func (h *TransactionHandler) HandledPayEvent(
	msg []byte,
	txStatus domain.TransactionStatus,
	sagaStep domain.SagaSteps,
	nextState domain.TransactionStatus,
	nextSagaState domain.SagaStatus,
	topic string,
	eventRepo repository.PgxEventRepo) {

	logger.Log.Info("Handling Account Balance Blocked event in TM", zap.ByteString("message", msg))

	// Idempotent processing
	var event domain.ConsumerEventForAccountService
	err := json.Unmarshal(msg, &event)
	if err != nil {
		logger.Log.Error("EVENT_UNMARSHAL_FAILED", zap.Error(err))
		return
	}

	processed, err := eventRepo.IsProcessed(event.EventID)
	if err != nil {
		logger.Log.Error("IDEMPOTENCY_CHECK_FAILED", zap.Error(err))
		return
	}
	if processed {
		logger.Log.Info("DUPLICATE_EVENT_IGNORED - Update Tx and SagaState", zap.String("event_id", event.EventID))
		return
	}

	logger.Log.Info("",
		zap.String("event_id", event.EventID),
		zap.String("event_type", event.EventType),
	)
	Tx, _, err := h.service.GetTransaction(event.TransactionID)
	if err != nil {
		logger.Log.Error("TRANSACTION_NOT_FOUND_FOR_EVENT",
			zap.String("event_id", event.EventID),
			zap.String("tx_id", event.TransactionID),
		)
		return
	}
	h.service.RecordSagaStep(Tx, string(sagaStep), domain.SagaStatusCompleted)
	h.service.UpdateTransactionWithSaga(Tx, txStatus, string(sagaStep)+"_"+domain.SagaStatusCompleted)
	if txStatus != domain.StatusCompleted {
		go h.event.PublishToAccountService(Tx, string(txStatus), topic, context.Background())
		go h.service.RecordSagaStep(Tx, string(nextSagaState), domain.SagaStatusInProgress)
		h.service.UpdateTransactionWithSaga(Tx, nextState, string(nextSagaState)+"_"+domain.SagaStatusInProgress)
	}

}
