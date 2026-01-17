package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	//"transaction-manager/internal/config"

	"transaction-manager/internal/config"
	"transaction-manager/internal/domain"
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
		Status:           "INITIATED",
		PaymentType:      req.PaymentType,
		PaymentMode:      req.PaymentMode,
		SagaStatus:       "STARTED",
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

	// emit event to account servicegit

	go h.event.PublishToAccountService(tx, "DEBIT_ACCOUNT", config.KAFKA_ACCOUNT_TOPIC, context.Background())

	// DEBIT_REQUESTED	REQUESTED
	h.service.RecordSagaStep(tx.ID, "DEBIT_REQUESTED", "REQUESTED")

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
