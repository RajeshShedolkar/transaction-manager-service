package events

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func NewReader(brokers []string, topic, group string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: group,
	})
}

func Consume(reader *kafka.Reader, handler func([]byte)) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka error:", err)
			continue
		}
		handler(m.Value)
	}
}
