package repository

import "transaction-manager/internal/domain"

type LedgerRepository interface {
	Append(entry *domain.LedgerEntry) error
	FindByTransactionID(txID string) ([]domain.LedgerEntry, error)

}