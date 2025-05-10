package user

import (
	"context"

	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

func (uc *UseCase) Get(ctx context.Context, id string) (*ports.UserDto, error) {
	uc.logger.Debug("Getting user", "user_id", id)

	userID, err := models.ParseUserID(id)
	if err != nil {
		uc.logger.Error("Failed to parse user ID", "error", err, "user_id", id)
		return nil, err
	}
	user, err := uc.userRepository.Get(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", id)
		return nil, err
	}

	dto := &ports.UserDto{
		ID:          user.ID().String(),
		Email:       user.Email(),
		Nickname:    user.Nickname(),
		Description: user.Description(),
		AvatarURL:   user.AvatarURL(),
	}

	uc.logger.Debug("User successfully retrieved", "user_id", id)

	return dto, nil
}
