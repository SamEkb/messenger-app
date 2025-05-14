package user

import (
	"context"

	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

func (uc *UseCase) GetByNickname(ctx context.Context, nickname string) (*ports.UserDto, error) {
	uc.logger.Debug("Getting user by nickname", "nickname", nickname)

	user, err := uc.userRepository.GetByNickname(ctx, nickname)
	if err != nil {
		uc.logger.Error("Failed to get user by nickname", "error", err, "nickname", nickname)
		return nil, err
	}

	dto := &ports.UserDto{
		ID:          user.ID().String(),
		Email:       user.Email(),
		Nickname:    user.Nickname(),
		Description: user.Description(),
		AvatarURL:   user.AvatarURL(),
	}

	uc.logger.Debug("User successfully retrieved", "nickname", nickname)

	return dto, nil
}
