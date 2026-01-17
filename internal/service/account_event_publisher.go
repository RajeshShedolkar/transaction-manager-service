package service

import (
	"context"
	"transaction-manager/internal/domain"
)

type AccountEventPublisher interface {
	PublishToAccountService(
		tx *domain.Transaction,
		eventType string,
		topic string,
		ctx context.Context,
	) error

	PublishToAccountServiceDLQ(
		event domain.TxEvent,
		topic string,
		errMsg string,
		ctx context.Context,
	) error
}
