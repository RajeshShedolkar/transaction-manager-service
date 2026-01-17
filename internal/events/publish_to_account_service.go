package events

// import (
// 	"context"
// 	"fmt"
// 	"transaction-manager/internal/config"
// 	"transaction-manager/internal/domain"
// )

// type KafkaAccountEventPublisher struct {
// }

// func NewKafkaAccountEventPublisher() *KafkaAccountEventPublisher {
// 	return &KafkaAccountEventPublisher{}
// }

// func (k *KafkaAccountEventPublisher) PublishToAccountService(tx *domain.Transaction, EventType string, topic string, ctx context.Context) error {

// 	cmd := domain.ConsumerEventForAccountService{
// 		TransactionID: tx.ID,
// 		EventType:     EventType,
// 		UserRefId:     tx.UserRefId,
// 		AccountRefId:  tx.SourceRefId,
// 		Amount:        tx.Amount,
// 	}
// 	accountEventWriter := NewKafkaProducer(config.KAFKA_BROKERS, topic)
// 	err := accountEventWriter.Publish(
// 		ctx,
// 		cmd.AccountRefId, // 🔑 partition key
// 		cmd,              // event payload
// 	)

// 	if err != nil {
// 		// Kafka did NOT ACK
// 		return fmt.Errorf("failed to publish debit command: %w", err)
// 	}

// 	return nil
// }

// func (k *KafkaAccountEventPublisher) PublishToDlqAccountService(event domain.TxEvent, topic string, err_msg string, ctx context.Context) error {

// 	cmd := domain.ConsumerDLQEventForCardService{
// 		SourceEventBody: event,
// 		ErrorMessage:    err_msg,
// 	}

// 	EventWriter := NewKafkaProducer(config.KAFKA_BROKERS, topic)
// 	err := EventWriter.Publish(
// 		ctx,
// 		"DLQ", // 🔑 partition key
// 		cmd,   // event payload
// 	)

// 	if err != nil {
// 		// Kafka did NOT ACK
// 		return fmt.Errorf("failed to publish debit command: %w", err)
// 	}

// 	return nil
// }
