package kfk

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

//
// ----------- CONSUMER -----------
//

func NewKafkaReader(brokers []string, topic, group string) *kafka.Reader {
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
			log.Println("Kafka consume error:", err)
			continue
		}

		handler(m.Value)
	}
}

//
// ----------- PRODUCER -----------
//

type Producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},    // ordering by key
		RequiredAcks: kafka.RequireAll, // wait for broker ACK
		Async:        false,            // synchronous produce
	}

	return &Producer{writer: writer}
}

func (p *Producer) Publish(
	ctx context.Context,
	key string,
	event any,
) error {

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
