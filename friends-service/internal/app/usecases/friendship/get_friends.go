package friendship

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
)

func (u *UseCase) GetFriends(ctx context.Context, userID string) ([]*ports.FriendshipDto, error) {
	u.logger.Info("getting friends")

	friends, err := u.friendRepository.GetFriends(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get friends", "error", err)
		return nil, err
	}

	result := make([]*ports.FriendshipDto, 0, len(friends))
	for _, f := range friends {
		dto, err := ports.NewFriendshipDto(
			f.ID().String(),
			f.RequestorID(),
			f.RecipientID(),
			string(f.Status()),
			f.CreatedAt(),
			f.UpdatedAt(),
		)
		if err != nil {
			u.logger.Error("failed to create Friendship DTO", "error", err)
			return nil, err
		}
		result = append(result, dto)
	}

	u.logger.Info("friends retrieved")
	return result, nil
}
