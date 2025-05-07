package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
)

type UserEventsKafkaProducer interface {
	ProduceUserRegisteredEvent(ctx context.Context, event *events.UserRegisteredEvent) error
}
