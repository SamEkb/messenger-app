package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/Shopify/sarama"
)

var _ ports.UserEventsKafkaProducer = (*UserEventsKafkaProducer)(nil)

type UserEventsKafkaProducer struct {
	producer sarama.SyncProducer
	logger   *log.Logger
	topic    string
}

func NewUserEventsKafkaProducer(cfg sarama.Config, kafkaCfg *env.KafkaConfig, logger *log.Logger) (*UserEventsKafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(kafkaCfg.Brokers, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &UserEventsKafkaProducer{
		producer: producer,
		logger:   logger,
		topic:    kafkaCfg.Topic,
	}, nil
}

func (p *UserEventsKafkaProducer) ProduceUserRegisteredEvent(ctx context.Context, event *events.UserRegisteredEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.GetUserId()),
		Value: sarama.ByteEncoder(data),
	}

	doneCh := make(chan struct{})
	var sendErr error
	var partition int32
	var offset int64

	go func() {
		defer close(doneCh)
		partition, offset, sendErr = p.producer.SendMessage(msg)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("message sending aborted: %w", ctx.Err())
	case <-doneCh:
		if sendErr != nil {
			return fmt.Errorf("failed to send message: %w", sendErr)
		}
		p.logger.Printf("Published UserRegisteredEvent for user %s to partition %d at offset %d",
			event.GetUserId(), partition, offset)
		return nil
	}
}

func (p *UserEventsKafkaProducer) Close() error {
	return p.producer.Close()
}
