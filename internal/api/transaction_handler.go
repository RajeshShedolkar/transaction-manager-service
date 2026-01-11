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
