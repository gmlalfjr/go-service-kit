package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

// KafkaService is a service for managing Kafka consumer
type KafkaService struct {
	reader   *kafka.Reader
	stopChan chan struct{}
	broker   string
}

// NewKafkaService creates a new KafkaService instance
func NewKafkaService(broker, topic, groupID string) *KafkaService {
	return &KafkaService{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
		stopChan: make(chan struct{}),
		broker:   broker,
	}
}

// Start starts the Kafka service
func (s *KafkaService) Start() error {
	log.Printf("[Kafka] Starting Kafka consumer on broker %s...", s.broker)
	go func() {
		for {
			select {
			case <-s.stopChan:
				log.Println("[Kafka] Kafka consumer stopped.")
				return
			default:
				msg, err := s.reader.ReadMessage(context.Background())
				if err != nil {
					log.Println("[Kafka] Error reading message:", err)
					continue
				}
				log.Printf("[Kafka] Message received: %s\n", string(msg.Value))
			}
		}
	}()
	return nil
}

// Stop stops the Kafka service
func (s *KafkaService) Stop() error {
	log.Println("[Kafka] Stopping Kafka consumer...")
	close(s.stopChan)
	return s.reader.Close()
}
