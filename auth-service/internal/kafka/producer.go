package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	events "github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/Shopify/sarama"
)

const (
	userEventsTopic = "user-events"
)

// Producer handles Kafka message production
type Producer struct {
	producer sarama.SyncProducer
}

// NewProducer creates a new Kafka producer
func NewProducer() (*Producer, error) {
	brokers := getBrokers()
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Producer{producer: producer}, nil
}

// Close closes the producer connection
func (p *Producer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}

// PublishUserRegistered publishes a UserRegisteredEvent to Kafka
func (p *Producer) PublishUserRegistered(_ context.Context, event *events.UserRegisteredEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: userEventsTopic,
		Key:   sarama.StringEncoder(event.GetUserId()),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Published UserRegisteredEvent for user %s to partition %d at offset %d",
		event.GetUserId(), partition, offset)
	return nil
}

// getBrokers gets the Kafka brokers from environment
func getBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// Default for local development
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}
