package user

import (
	"context"

	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

func (uc *UseCase) Create(ctx context.Context, dto *ports.UserDto) (string, error) {
	uc.logger.Debug("Creating user", "username", dto.Nickname, "email", dto.Email)

	userID, err := models.ParseUserID(dto.ID)
	if err != nil {
		uc.logger.Error("Failed to parse user ID", "error", err, "user_id", dto.ID)
		return "", err
	}
	newUser, err := models.NewUser(userID, dto.Email, dto.Nickname, dto.Description, dto.AvatarURL)
	if err != nil {
		uc.logger.Error("Failed to create user", "error", err, "user_id", dto.ID)
		return "", err
	}
	userID, err = uc.userRepository.Create(ctx, newUser)
	if err != nil {
		uc.logger.Error("Failed to create user", "error", err, "user_id", dto.ID)
		return "", err
	}

	return userID.String(), nil
}
