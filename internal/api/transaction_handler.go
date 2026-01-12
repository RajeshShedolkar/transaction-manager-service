package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"transaction-manager/internal/domain"
	"transaction-manager/internal/service"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tx := &domain.Transaction{
		PaymentType: req.PaymentType,
		PaymentMode: req.PaymentMode,
		Amount:      req.Amount,
		Currency:    req.Currency,
	}

	err := h.service.CreateImmediateTransaction(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := CreateTransactionResponse{
		TransactionID: tx.ID,
		Status:        string(tx.Status),
		Message:       "Transaction processed successfully",
	}

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

