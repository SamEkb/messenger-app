package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (u *UseCase) CreateChat(ctx context.Context, participants []string) (*ports.ChatDto, error) {
	u.logger.Info("creating chat")

	if len(participants) < 2 {
		return nil, errors.NewInvalidInputError("chat requires at least 2 participants")
	}

	profilesResp, err := u.userClient.GetProfiles(ctx, &users.GetProfilesRequest{
		UserIds: participants,
	})
	if err != nil {
		u.logger.Error("failed to get user profiles", "error", err)
		return nil, fmt.Errorf("failed to get user profiles: %w", err)
	}

	if len(profilesResp.NotFoundIds) > 0 {
		notFoundUsers := strings.Join(profilesResp.NotFoundIds, ", ")
		u.logger.Info("some users not found", "missing_users", notFoundUsers)
		return nil, errors.NewNotFoundError(fmt.Sprintf("users not found: %s", notFoundUsers))
	}

	friendshipResp, err := u.friendClient.CheckFriendshipsStatus(ctx, &friends.CheckFriendshipsStatusRequest{
		UserIds: participants,
	})
	if err != nil {
		u.logger.Error("failed to check friendships status", "error", err)
		return nil, err
	}

	if !friendshipResp.AllAreFriends {
		if len(friendshipResp.NonFriendPairs) > 0 {
			pair := friendshipResp.NonFriendPairs[0]
			return nil, errors.NewForbiddenError(
				fmt.Sprintf("users %s and %s are not friends", pair.UserID1, pair.UserID2))
		}
		return nil, errors.NewForbiddenError("some participants are not friends")
	}

	var chat *models.Chat
	err = u.txManager.RunTx(ctx, func(sessionCtx mongo.SessionContext) error {
		var err error
		chat, err = u.chatRepository.Create(sessionCtx, participants)
		if err != nil {
			u.logger.Error("failed to create chat", "error", err)
			return fmt.Errorf("failed to create chat: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ports.NewChatDto(
		chat.ID().String(),
		chat.Participants(),
		make([]*ports.MessageDto, 0),
		chat.CreatedAt(),
		chat.UpdatedAt(),
	), nil
}
