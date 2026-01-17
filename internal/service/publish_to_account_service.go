package service

import (
	"context"
	"transaction-manager/internal/domain"
)

type KafkaAccountEventPublisherImpl struct {
}

func NewKafkaAccountEventPublisher() *KafkaAccountEventPublisherImpl {
	return &KafkaAccountEventPublisherImpl{}
}


func (k *KafkaAccountEventPublisherImpl) PublishToAccountService(tx *domain.Transaction, EventType string, topic string, ctx context.Context) error {

	// cmd := domain.ConsumerEventForAccountService{
	// 	TransactionID: tx.ID,
	// 	EventType:     EventType,
	// 	UserRefId:     tx.UserRefId,
	// 	AccountRefId:  tx.SourceRefId,
	// 	Amount:        tx.Amount,
	// }
	// accountEventWriter := NewKafkaProducer(config.KAFKA_BROKERS, topic)
	// err := accountEventWriter.Publish(
	// 	ctx,
	// 	cmd.AccountRefId, // 🔑 partition key
	// 	cmd,              // event payload
	// )

	// if err != nil {
	// 	// Kafka did NOT ACK
	// 	return fmt.Errorf("failed to publish debit command: %w", err)
	// }

	return nil
}

func (k *KafkaAccountEventPublisherImpl) PublishToAccountServiceDLQ(event domain.TxEvent, topic string, err_msg string, ctx context.Context) error {

	// cmd := domain.ConsumerDLQEventForCardService{
	// 	SourceEventBody: event,
	// 	ErrorMessage:    err_msg,
	// }

	// EventWriter := NewKafkaProducer(config.KAFKA_BROKERS, topic)
	// err := EventWriter.Publish(
	// 	ctx,
	// 	"DLQ", // 🔑 partition key
	// 	cmd,   // event payload
	// )

	// if err != nil {
	// 	// Kafka did NOT ACK
	// 	return fmt.Errorf("failed to publish debit command: %w", err)
	// }

	return nil
}
