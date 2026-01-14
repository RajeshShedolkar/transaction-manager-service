package repository

import (
	"context"
	"transaction-manager/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxTransactionRepo struct {
	db *pgxpool.Pool
}

func NewPgxTransactionRepo(db *pgxpool.Pool) *PgxTransactionRepo {
	return &PgxTransactionRepo{db: db}
}

func (r *PgxTransactionRepo) Save(tx *domain.Transaction) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO transactions 
	(id, user_ref_id, source_ref_id, destination_ref_id, payment_type, payment_mode, status, dc_flag, amount, currency, network_txn_id, gateway_txn_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		tx.ID,
		tx.UserRefId,
		tx.SourceRefId,
		tx.DestinationRefId,
		tx.PaymentType,
		tx.PaymentMode,
		tx.Status,
		tx.DcFlag,
		tx.Amount,
		tx.Currency,
		tx.NetworkTxnId,
		tx.GatewayTxnId,
	)

	return err
}

func (r *PgxTransactionRepo) FindByID(id string) (*domain.Transaction, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id,payment_type,payment_mode,status,amount,currency 
		 FROM transactions WHERE id=$1`, id)

	var tx domain.Transaction
	err := row.Scan(&tx.ID, &tx.PaymentType, &tx.PaymentMode, &tx.Status, &tx.Amount, &tx.Currency)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *PgxTransactionRepo) UpdateStatus(id string, status domain.TransactionStatus) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE transactions SET status=$1, updated_at=NOW() WHERE id=$2`,
		status, id,
	)
	return err
}


func (r *PgxTransactionRepo) FindByNetworkTxnID(id string) (*domain.Transaction, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id,payment_type,payment_mode,status,amount,currency 
		 FROM transactions WHERE network_txn_id=$1`, id)

	var tx domain.Transaction
	err := row.Scan(&tx.ID, &tx.PaymentType, &tx.PaymentMode, &tx.Status, &tx.Amount, &tx.Currency)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}