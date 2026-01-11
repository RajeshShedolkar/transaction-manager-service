package service

import "transaction-manager/internal/domain"

type TransactionService interface {
	CreateImmediateTransaction(tx *domain.Transaction) error
}
