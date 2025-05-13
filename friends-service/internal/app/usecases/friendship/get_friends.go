package friendship

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
)

func (u *UseCase) GetFriends(ctx context.Context, userID string) ([]*ports.FriendshipDto, error) {
	u.logger.Info("getting friends")

	friends, err := u.friendRepository.GetFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	friendIDs := make([]string, 0, len(friends))
	for _, f := range friends {
		friendIDs = append(friendIDs, f.RecipientID())
	}

	request := &users.GetProfilesRequest{UserIds: friendIDs}
	profilesResp, err := u.userClient.GetProfiles(ctx, request)
	if err != nil {
		return nil, err
	}

	result := make([]*ports.FriendshipDto, 0, len(friends))
	for _, f := range friends {
		profile := profilesResp.Profiles[f.RecipientID()]
		dto, _ := ports.NewFriendshipDto(
			f.ID().String(),
			f.RequestorID(),
			f.RecipientID(),
			profile.Nickname,
			profile.AvatarURL,
			string(f.Status()),
			f.CreatedAt(),
			f.UpdatedAt(),
		)
		result = append(result, dto)
	}

	u.logger.Info("friends retrieved")
	return result, nil
}
