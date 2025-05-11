package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/SamEkb/messenger-app/users-service/config/env"
	"github.com/Shopify/sarama"
)

const (
	userEventsTopic = "user-events"
)

type EventHandler interface {
	HandleUserRegistered(ctx context.Context, event *events.UserRegisteredEvent) error
}

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       EventHandler
	topics        []string
	ready         chan bool
	canceller     func()
}

func NewConsumerWithConfig(handler EventHandler, kafkaConfig *env.KafkaConfig) (*Consumer, error) {
	if kafkaConfig == nil {
		return nil, fmt.Errorf("kafka config is nil")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(kafkaConfig.Brokers, kafkaConfig.ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		topics:        []string{kafkaConfig.Topic},
		ready:         make(chan bool),
	}, nil
}

type consumerGroupHandler struct {
	handler EventHandler
	ready   chan bool
}

func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group session setup: member ID = %s", session.MemberID())
	close(h.ready)
	return nil
}

func (h *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group session cleanup: member ID = %s", session.MemberID())
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received message from topic %s, partition %d, offset %d",
			msg.Topic, msg.Partition, msg.Offset)

		if msg.Topic == userEventsTopic {
			var event events.UserRegisteredEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to unmarshal user event: %v", err)
				continue
			}

			if err := h.handler.HandleUserRegistered(session.Context(), &event); err != nil {
				log.Printf("Error handling message: %v", err)
				continue
			}
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.canceller = cancel

	handler := &consumerGroupHandler{
		handler: c.handler,
		ready:   c.ready,
	}

	go func() {
		for {
			err := c.consumerGroup.Consume(ctx, c.topics, handler)
			if err != nil {
				log.Printf("Error from consumer group: %v", err)
			}

			if ctx.Err() != nil {
				log.Println("Context cancelled, stopping consumer")
				return
			}

			c.ready = make(chan bool)
			handler.ready = c.ready
		}
	}()

	<-c.ready
	log.Println("Kafka consumer started")
	return nil
}

func (c *Consumer) Close() error {
	if c.canceller != nil {
		c.canceller()
	}
	if c.consumerGroup != nil {
		return c.consumerGroup.Close()
	}
	return nil
}

func (c *Consumer) Ready() <-chan bool {
	return c.ready
}

func getBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// Default for local development
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}
