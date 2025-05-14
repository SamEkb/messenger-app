package kafka

import (
	"context"
	"encoding/json"

	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/Shopify/sarama"
)

var _ ports.UserEventsKafkaProducer = (*UserEventsKafkaProducer)(nil)

type UserEventsKafkaProducer struct {
	producer sarama.SyncProducer
	logger   logger.Logger
	topic    string
}

func NewUserEventsKafkaProducer(kafkaCfg *env.KafkaConfig, logger logger.Logger) (*UserEventsKafkaProducer, error) {
	if kafkaCfg == nil {
		return nil, errors.NewInvalidInputError("kafka config is nil")
	}

	cfg := NewSaramaConfig(kafkaCfg, logger)

	producer, err := sarama.NewSyncProducer(kafkaCfg.Brokers, cfg)
	if err != nil {
		return nil, errors.NewServiceError(err, "failed to create Kafka producer")
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
		return errors.NewInternalError(err, "failed to marshal event")
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
		return errors.NewTimeoutError("message sending aborted: %v", ctx.Err()).
			WithDetails("user_id", event.GetUserId())
	case <-doneCh:
		if sendErr != nil {
			p.logger.Error("failed to send message", "error", sendErr)
			return errors.NewServiceError(sendErr, "failed to send message").
				WithDetails("user_id", event.GetUserId())
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
	err := p.producer.Close()
	if err != nil {
		return errors.NewServiceError(err, "failed to close Kafka producer")
	}
	return nil
}
