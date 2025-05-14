package auth

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *UseCase) Register(ctx context.Context, dto *ports.RegisterDto) (models.UserID, error) {
	a.logger.Debug("register attempt", "username", dto.Username, "email", dto.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("failed to hash password", "error", err)
		return models.UserID{}, errors.NewInternalError(err, "failed to hash password")
	}

	userID := models.UserID(uuid.New())

	user, err := models.NewUser(
		userID,
		dto.Username,
		dto.Email,
		hashedPassword,
	)
	if err != nil {
		a.logger.Error("failed to create user model", "error", err)
		return models.UserID{}, err
	}

	if err = a.authRepo.Create(ctx, user); err != nil {
		a.logger.Error("failed to save user", "error", err)
		return models.UserID{}, err
	}

	event := &events.UserRegisteredEvent{
		UserId:       user.ID().String(),
		Username:     user.Username(),
		Email:        user.Email(),
		RegisteredAt: timestamppb.Now(),
	}

	if err = a.userEventPublisher.ProduceUserRegisteredEvent(ctx, event); err != nil {
		a.logger.Warn("failed to publish user registration event", "error", err, "user_id", user.ID())
		return models.UserID{}, err
	}

	a.logger.Info("user registered successfully", "user_id", user.ID(), "username", user.Username(), "email", user.Email())
	return user.ID(), nil
}
