package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/Shopify/sarama"
)

const (
	userEventsTopic = "user-events"
)

type EventHandler interface {
	HandleUserRegistered(ctx context.Context, event *events.UserRegisteredEvent) error
}

type Consumer struct {
	consumer  sarama.Consumer
	handler   EventHandler
	topics    []string
	ready     chan bool
	canceller func()
}

func NewConsumer(handler EventHandler) (*Consumer, error) {
	brokers := getBrokers()
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		handler:  handler,
		topics:   []string{userEventsTopic},
		ready:    make(chan bool),
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.canceller = cancel

	for _, topic := range c.topics {
		partitions, err := c.consumer.Partitions(topic)
		if err != nil {
			return fmt.Errorf("failed to get partitions for topic %s: %w", topic, err)
		}

		for _, partition := range partitions {
			pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				return fmt.Errorf("failed to start consumer for topic %s partition %d: %w", topic, partition, err)
			}

			go func(pc sarama.PartitionConsumer) {
				defer pc.Close()
				for {
					select {
					case msg := <-pc.Messages():
						if err := c.handleMessage(ctx, msg); err != nil {
							log.Printf("Error handling message: %v", err)
						}
					case err := <-pc.Errors():
						log.Printf("Consumer error: %v", err)
					case <-ctx.Done():
						return
					}
				}
			}(pc)
		}
	}

	log.Println("Kafka consumer started")
	c.ready <- true
	return nil
}

func (c *Consumer) Close() error {
	if c.canceller != nil {
		c.canceller()
	}
	if c.consumer != nil {
		return c.consumer.Close()
	}
	return nil
}

func (c *Consumer) Ready() <-chan bool {
	return c.ready
}

func (c *Consumer) handleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	log.Printf("Received message from topic %s, partition %d, offset %d",
		msg.Topic, msg.Partition, msg.Offset)

	if msg.Topic == userEventsTopic {
		var event events.UserRegisteredEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user event: %w", err)
		}

		return c.handler.HandleUserRegistered(ctx, &event)
	}

	return fmt.Errorf("unknown topic: %s", msg.Topic)
}

func getBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// Default for local development
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}
