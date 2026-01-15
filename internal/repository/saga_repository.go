package repository

import "transaction-manager/internal/domain"

type SagaRepository interface {
	AddStep(step *domain.SagaStep) error
	UpdateSagaStatus(txID string, status string) error
}
