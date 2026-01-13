package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxEventRepo struct {
	db *pgxpool.Pool
}

func NewPgxEventRepo(db *pgxpool.Pool) *PgxEventRepo {
	return &PgxEventRepo{db: db}
}

func (r *PgxEventRepo) IsProcessed(eventID string) (bool, error) {

	row := r.db.QueryRow(context.Background(),
		`SELECT event_id FROM processed_events WHERE event_id=$1`, eventID)

	var id string
	err := row.Scan(&id)

	if err != nil {
		return false, nil // not found
	}
	return true, nil
}

func (r *PgxEventRepo) MarkProcessed(eventID string, eventType string) error {

	_, err := r.db.Exec(context.Background(),
		`INSERT INTO processed_events (event_id, event_type)
		 VALUES ($1,$2)`,
		eventID, eventType,
	)
	return err
}
