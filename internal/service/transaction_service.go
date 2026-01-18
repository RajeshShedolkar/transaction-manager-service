package service

import (
	"transaction-manager/internal/domain"
)

type TransactionService interface {
	CreateImmediateTransaction(tx *domain.Transaction) error
	GetTransaction(id string) (*domain.Transaction, []domain.LedgerEntry, error)
	GetTransactionLedger(id string) (*domain.Transaction, []domain.LedgerEntry, error)
	CreateNEFTTransaction(tx *domain.Transaction) error
	HandleNEFTSettlement(txID string, status string) error
	RecordSagaStep(tx *domain.Transaction, step, status string)
	UpdateSagaStatus(txID, status string)
	UpdateTransactionWithSaga(tx *domain.Transaction, status domain.TransactionStatus, sagaCurrState string) error
}
