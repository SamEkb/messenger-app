package user

import (
	"context"

	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

func (uc *UseCase) Create(ctx context.Context, dto *ports.UserDto) (string, error) {
	uc.logger.Debug("Creating user", "username", dto.Nickname, "email", dto.Email)

	parsedID, err := models.ParseUserID(dto.ID)
	if err != nil {
		uc.logger.Error("Failed to parse user ID", "error", err, "user_id", dto.ID)
		return "", err
	}

	newUser, err := models.NewUser(parsedID, dto.Email, dto.Nickname, dto.Description, dto.AvatarURL)
	if err != nil {
		uc.logger.Error("Failed to create user", "error", err, "user_id", dto.ID)
		return "", err
	}

	var userID models.UserID
	err = uc.txManager.RunTx(ctx, func(txCtx context.Context) error {
		var err error
		userID, err = uc.userRepository.Create(txCtx, newUser)
		if err != nil {
			uc.logger.Error("Failed to create user", "error", err, "user_id", dto.ID)
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return userID.String(), nil
}
