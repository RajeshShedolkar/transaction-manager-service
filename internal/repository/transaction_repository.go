package repository

import "transaction-manager/internal/domain"

type TransactionRepository interface {
	Save(tx *domain.Transaction) error
	FindByID(id string) (*domain.Transaction, error)
	FindByNetworkTxnID(id string) (*domain.Transaction, error)
	UpdateStatus(id string, status domain.TransactionStatus) error
	UpdateStatusWithSaga(txID string, status domain.TransactionStatus, sagaState string) error
}
