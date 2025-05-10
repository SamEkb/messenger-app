package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/models"
)

type FriendshipRepository interface {
	Create(ctx context.Context, friendship *models.Friendship) error
	GetByUserIDs(ctx context.Context, userID1, userID2 string) (*models.Friendship, error)
	GetAllFriendships(ctx context.Context, userID string) ([]*models.Friendship, error)
	Update(ctx context.Context, friendship *models.Friendship) error
	Delete(ctx context.Context, friendship *models.Friendship) error
	CheckFriendships(ctx context.Context, userIDs []string) ([]UserPair, bool, error)
}

type UserPair struct {
	UserID1 string
	UserID2 string
}
