package service

import (
	"transaction-manager/internal/domain"
)

type TransactionService interface {
	CreateImmediateTransaction(tx *domain.Transaction) error
	GetTransaction(id string) (*domain.Transaction, []domain.LedgerEntry, error)
	CreateNEFTTransaction(tx *domain.Transaction) error
	HandleNEFTSettlement(txID string, status string) error
	RecordSagaStep(txID, step, status string)
	UpdateSagaStatus(txID, status string)
	UpdateTransactionWithSaga(txID *domain.Transaction, status domain.TransactionStatus, sagaStatus domain.SagaStatus) error
}
