package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/Shopify/sarama"
)

var _ ports.UserEventsKafkaProducer = (*UserEventsKafkaProducer)(nil)

type UserEventsKafkaProducer struct {
	producer sarama.SyncProducer
	logger   *slog.Logger
	topic    string
}

func NewUserEventsKafkaProducer(cfg *sarama.Config, kafkaCfg *env.KafkaConfig, logger *slog.Logger) (*UserEventsKafkaProducer, error) {
	if cfg == nil {
		panic("kafka config is nil")
	}

	producer, err := sarama.NewSyncProducer(kafkaCfg.Brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &UserEventsKafkaProducer{
		producer: producer,
		logger:   logger.With("component", "kafka_producer"),
		topic:    kafkaCfg.Topic,
	}, nil
}

func (p *UserEventsKafkaProducer) ProduceUserRegisteredEvent(ctx context.Context, event *events.UserRegisteredEvent) error {
	p.logger.Debug("preparing to produce user registered event", "user_id", event.GetUserId())

	data, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("failed to marshal event", "error", err)
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
		p.logger.Warn("message sending aborted", "error", ctx.Err())
		return fmt.Errorf("message sending aborted: %w", ctx.Err())
	case <-doneCh:
		if sendErr != nil {
			p.logger.Error("failed to send message", "error", sendErr)
			return fmt.Errorf("failed to send message: %w", sendErr)
		}
		p.logger.Info("published UserRegisteredEvent",
			"user_id", event.GetUserId(),
			"partition", partition,
			"offset", offset,
			"topic", p.topic)
		return nil
	}
}

func (p *UserEventsKafkaProducer) Close() error {
	p.logger.Info("closing kafka producer")
	return p.producer.Close()
}
