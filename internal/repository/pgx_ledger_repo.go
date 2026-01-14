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
	(id, transaction_id, account_ref_id, dc_flag, entry_type, amount, source)
	VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		entry.ID,
		entry.TransactionID,
		entry.AccountRefId,
		entry.DcFlag,
		entry.EntryType,
		entry.Amount,
		entry.Source,
	)

	return err
}

func (r *PgxLedgerRepo) FindByTransactionID(txID string) ([]domain.LedgerEntry, error) {

	rows, err := r.db.Query(context.Background(),
		`SELECT id,transaction_id,entry_type,amount,source 
		 FROM ledger_entries WHERE transaction_id=$1`, txID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.LedgerEntry
	for rows.Next() {
		var e domain.LedgerEntry
		err := rows.Scan(&e.ID, &e.TransactionID, &e.EntryType, &e.Amount, &e.Source)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
