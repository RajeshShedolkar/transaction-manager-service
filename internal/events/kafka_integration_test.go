package events

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "localhost:9092"
	testTopic   = "test-topic1"
)

func TestKafkaProduceConsume(t *testing.T) {

	// ---- Producer ----
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
		Topic:   testTopic,
	})
	defer writer.Close()

	message := kafka.Message{
		Key:   []byte("test-key1"),
		Value: []byte(`{"eventId":"TEST1","eventType":"CARD_AUTH","amount":500}`),
	}

	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		t.Fatalf("failed to write kafka message: %v", err)
	}

	log.Println("Kafka message produced")

	// ---- Consumer ----
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   testTopic,
		GroupID: "tm-test-group",
	})
	defer reader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		t.Fatalf("failed to read kafka message: %v", err)
	}

	log.Println("Kafka message consumed:", string(msg.Value))

	// ---- Assertion ----
	if string(msg.Value) == "" {
		t.Fatal("empty kafka message received")
	}

	if string(msg.Value) != string(message.Value) {
		t.Fatalf("message mismatch. expected %s, got %s",
			string(message.Value), string(msg.Value))
	}
}
