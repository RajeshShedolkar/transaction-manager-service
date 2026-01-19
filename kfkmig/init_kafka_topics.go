package kfkmig

import (
	"transaction-manager/internal/config"

	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type KafkaAdmin struct {
	brokers []string
	config  *sarama.Config
}

// Constructor
func NewKafkaAdmin(brokers []string) (*KafkaAdmin, error) {
	cfg := sarama.NewConfig()

	// Kafka version (important)
	cfg.Version = sarama.V2_8_0_0

	// Admin timeouts
	cfg.Admin.Timeout = 10 * time.Second

	return &KafkaAdmin{
		brokers: brokers,
		config:  cfg,
	}, nil
}

// EnsureTopic ensures topic exists, otherwise creates it
func (k *KafkaAdmin) EnsureTopic(
	topic string,
	partitions int32,
	replicationFactor int16,
) error {

	admin, err := sarama.NewClusterAdmin(k.brokers, k.config)
	if err != nil {
		return fmt.Errorf("create cluster admin failed: %w", err)
	}
	defer admin.Close()

	// Fetch existing topics
	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("list topics failed: %w", err)
	}

	// Topic already exists
	if _, exists := topics[topic]; exists {
		return nil
	}

	// Topic does NOT exist → create
	detail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries: map[string]*string{
			"cleanup.policy":      stringPtr("delete"),
			"retention.ms":        stringPtr("604800000"), // 7 days
			"segment.ms":          stringPtr("86400000"),  // 1 day
			"min.insync.replicas": stringPtr("1"),
		},
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = admin.CreateTopic(topic, detail, false)
	if err != nil {
		return fmt.Errorf("create topic %s failed: %w", topic, err)
	}

	fmt.Println("Process of Kafka Topic Creation Completed...")

	return nil
}

// Helper
func stringPtr(v string) *string {
	return &v
}

func BootstrapKafkaTopics() error {
	admin, err := NewKafkaAdmin([]string{"kafka:9092"})
	if err != nil {
		return err
	}
	var AllPaymentTopics = append(config.AllIMPSTopics, config.AllNEFTTopics...)

	for _, topic := range AllPaymentTopics {
		err := admin.EnsureTopic(topic, 3, 1)
		if err != nil {
			return err
		}
	}

	return nil
}

// func InitKafkaTopics() {
// 	brokers := config.KAFKA_BROKERS
// 	acc_err := CreateTopicIfNotExists(brokers[0], config.KAFKA_ACCOUNT_TOPIC)
// 	if acc_err != nil {
// 		logger.Log.Fatal("Failed to create Kafka topic", zap.Error(acc_err))
// 	} else {
// 		logger.Log.Info("Kafka topic ensured", zap.String("topic", config.KAFKA_ACCOUNT_TOPIC))
// 	}

// 	dlq_err := CreateTopicIfNotExists(brokers[0], config.KAFKA_DLQ_TOPIC)
// 	if dlq_err != nil {
// 		logger.Log.Fatal("Failed to create Kafka topic", zap.Error(dlq_err))
// 	} else {
// 		logger.Log.Info("Kafka topic ensured", zap.String("topic", config.KAFKA_DLQ_TOPIC))
// 	}

// 	card_err := CreateTopicIfNotExists(brokers[0], config.KAFKA_CARD_EVENT_TOPIC)
// 	if card_err != nil {
// 		logger.Log.Fatal("Failed to create Kafka topic", zap.Error(card_err))
// 	} else {
// 		logger.Log.Info("Kafka topic ensured", zap.String("topic", config.KAFKA_CARD_EVENT_TOPIC))
// 	}
// }

func main(){
	
}