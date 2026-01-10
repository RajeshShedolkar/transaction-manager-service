package repository

import "transaction-manager/internal/domain"

type LedgerRepository interface {
	Append(entry *domain.LedgerEntry) error
}
