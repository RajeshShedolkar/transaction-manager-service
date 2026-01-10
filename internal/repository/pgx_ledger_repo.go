package repository

import (
	"context"
	"transaction-manager/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxLedgerRepo struct {
	db *pgxpool.Pool
}

func NewPgxLedgerRepo(db *pgxpool.Pool) *PgxLedgerRepo {
	return &PgxLedgerRepo{db: db}
}

func (r *PgxLedgerRepo) Append(entry *domain.LedgerEntry) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO ledger_entries 
		(id, transaction_id, entry_type, amount, source)
		VALUES ($1,$2,$3,$4,$5)`,
		entry.ID, entry.TransactionID, entry.EntryType, entry.Amount, entry.Source,
	)
	return err
}
