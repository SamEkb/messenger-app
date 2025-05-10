package kafka

import (
	"context"

	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

func (s *UsersServiceServer) HandleUserRegistered(ctx context.Context, event *events.UserRegisteredEvent) error {
	s.logger.Info("Handling user registered event",
		"user_id", event.UserId,
		"username", event.Username,
		"email", event.Email)

	dto := &ports.UserDto{
		ID:          event.UserId,
		Email:       event.Email,
		Nickname:    event.Username,
		Description: "",
		AvatarURL:   "",
	}

	userID, err := s.userUseCase.Create(ctx, dto)
	if err != nil {
		s.logger.Error("Failed to create user from event",
			"error", err,
			"user_id", event.UserId)
		return err
	}

	s.logger.Info("User successfully created from Kafka event",
		"user_id", userID)
	return nil
}
