package user

import (
	"context"

	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

func (uc *UseCase) Update(ctx context.Context, dto *ports.UserDto) error {
	uc.logger.Debug("Updating user", "user_id", dto.ID)

	id, err := models.ParseUserID(dto.ID)
	if err != nil {
		uc.logger.Error("Failed to parse user ID", "error", err, "user_id", dto.ID)
		return err
	}

	user, err := models.NewUser(
		id,
		dto.Email,
		dto.Nickname,
		dto.Description,
		dto.AvatarURL,
	)
	if err != nil {
		uc.logger.Error("Failed to update user", "error", err, "user_id", dto.ID)
		return err
	}

	if err := uc.userRepository.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user", "error", err, "user_id", dto.ID)
		return err
	}

	uc.logger.Debug("User successfully updated", "user_id", dto.ID)

	return nil
}
