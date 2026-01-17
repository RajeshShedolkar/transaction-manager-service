package events

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

func CreateTopicIfNotExists(broker string, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+fmt.Sprint(controller.Port))
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	topics := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     6,
			ReplicationFactor: 1, // change to 3 in prod
		},
	}

	err = controllerConn.CreateTopics(topics...)
	if err != nil {
		// Topic already exists → ignore
		if err.Error() == "Topic with this name already exists" {
			return nil
		}
		return err
	}

	return nil
}
