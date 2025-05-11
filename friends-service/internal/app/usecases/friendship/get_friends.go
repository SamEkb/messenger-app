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
		u.logger.Error("failed to get friends", "error", err)
		return nil, err
	}

	friendIDs := make([]string, 0, len(friends))
	for _, f := range friends {
		friendIDs = append(friendIDs, f.RecipientID())
	}

	request := &users.GetProfilesRequest{UserIds: friendIDs}
	profilesResp, err := u.userClient.GetProfiles(ctx, request)
	if err != nil {
		u.logger.Error("failed to get user profiles", "error", err)
		return nil, err
	}

	profileMap := make(map[string]*ports.UserProfile, len(profilesResp.Profiles))
	for _, p := range profilesResp.Profiles {
		profileMap[p.UserID] = p
	}

	result := make([]*ports.FriendshipDto, 0, len(friends))
	for _, f := range friends {
		profile := profileMap[f.RecipientID()]
		nickname := ""
		avatarURL := ""
		if profile != nil {
			nickname = profile.Nickname
			avatarURL = profile.AvatarURL
		}
		dto, err := ports.NewFriendshipDto(
			f.ID().String(),
			f.RequestorID(),
			f.RecipientID(),
			nickname,
			avatarURL,
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
