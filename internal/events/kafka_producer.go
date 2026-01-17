package events

import (
	"github.com/segmentio/kafka-go"
)

var KafkaCommandWriter *kafka.Writer

func InitProducer(topic string, brokers []string) {
	KafkaCommandWriter = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{}, // important for ordering
		RequiredAcks: kafka.RequireAll,
		Async:        false, // we want ACK
	}
}
