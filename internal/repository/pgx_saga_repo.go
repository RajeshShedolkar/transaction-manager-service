package repository

import (
	"context"
	"transaction-manager/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxSagaRepo struct {
	db *pgxpool.Pool
}

func NewPgxSagaRepo(db *pgxpool.Pool) *PgxSagaRepo {
	return &PgxSagaRepo{db: db}
}

func (r *PgxSagaRepo) AddStep(step *domain.SagaStep) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO saga_steps
		 (id, transaction_id, step_name, status, tx_state)
		 VALUES ($1,$2,$3,$4, $5)`,
		step.ID,
		step.TransactionID,
		step.StepName,
		step.Status,
		step.TxState,
	)
	return err
}

func (r *PgxSagaRepo) UpdateSagaStatus(txID string, status string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE transactions SET saga_status=$1 WHERE id=$2`,
		status,
		txID,
	)
	return err
}
