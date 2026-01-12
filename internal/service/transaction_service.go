package service

import "transaction-manager/internal/domain"

type TransactionService interface {
	CreateImmediateTransaction(tx *domain.Transaction) error
	GetTransaction(id string) (*domain.Transaction, []domain.LedgerEntry, error)

}
