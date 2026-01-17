package events

import (
	"transaction-manager/internal/config"
	"transaction-manager/pkg/logger"

	"go.uber.org/zap"
)

func InitKafkaTopics() {
	brokers := config.KAFKA_BROKERS
	acc_err := CreateTopicIfNotExists(brokers[0], config.KAFKA_ACCOUNT_TOPIC)
	if acc_err != nil {
		logger.Log.Fatal("Failed to create Kafka topic", zap.Error(acc_err))
	} else {
		logger.Log.Info("Kafka topic ensured", zap.String("topic", config.KAFKA_ACCOUNT_TOPIC))
	}

	dlq_err := CreateTopicIfNotExists(brokers[0], config.KAFKA_DLQ_TOPIC)
	if dlq_err != nil {
		logger.Log.Fatal("Failed to create Kafka topic", zap.Error(dlq_err))
	} else {
		logger.Log.Info("Kafka topic ensured", zap.String("topic", config.KAFKA_DLQ_TOPIC))
	}

	card_err := CreateTopicIfNotExists(brokers[0], config.KAFKA_CARD_EVENT_TOPIC)
	if card_err != nil {
		logger.Log.Fatal("Failed to create Kafka topic", zap.Error(card_err))
	} else {
		logger.Log.Info("Kafka topic ensured", zap.String("topic", config.KAFKA_CARD_EVENT_TOPIC))
	}
}
