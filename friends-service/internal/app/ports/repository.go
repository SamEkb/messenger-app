package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/models"
)

type FriendshipRepository interface {
	GetFriends(ctx context.Context, userID string) ([]*models.Friendship, error)
	SendFriendRequest(ctx context.Context, requestorID, recipientID string) error
	AcceptFriendRequest(ctx context.Context, recipientID, requestorID string) error
	RejectFriendRequest(ctx context.Context, recipientID, requestorID string) error
	Delete(ctx context.Context, userID string, friendID string) error
}
