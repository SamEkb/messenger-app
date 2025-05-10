package in_memory

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/models"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/friends-service/pkg/errors"
)

var _ ports.FriendshipRepository = (*FriendshipRepository)(nil)

type FriendshipRepository struct {
	friendships map[string]*models.Friendship
	mx          sync.RWMutex
	logger      *slog.Logger
}

func NewFriendshipRepository(logger *slog.Logger) *FriendshipRepository {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	return &FriendshipRepository{
		friendships: make(map[string]*models.Friendship),
	}
}

func (r *FriendshipRepository) Create(ctx context.Context, friendship *models.Friendship) error {
	r.logger.Info("creating friendship",
		slog.String("user_id", friendship.UserID()), slog.String("friend_id", friendship.FriendID()),
	)
	r.mx.Lock()
	defer r.mx.Unlock()

	key := fmt.Sprintf("%s:%s", friendship.UserID(), friendship.FriendID())
	if _, exists := r.friendships[key]; exists {
		return errors.NewAlreadyExistsError("friendship already exists", friendship.UserID(), friendship.FriendID())
	}
	r.friendships[key] = friendship

	r.logger.Info("friendship created",
		slog.String("user_id", friendship.UserID()), slog.String("friend_id", friendship.FriendID()),
	)
	return nil

}

func (r *FriendshipRepository) GetByUserIDs(ctx context.Context, userID1, userID2 string) (*models.Friendship, error) {
	if userID1 == userID2 {
		return nil, errors.NewInvalidInputError("user IDs must be different",
			"userID1", userID1, "userID2", userID2)
	}

	key1 := fmt.Sprintf("%s:%s", userID1, userID2)
	key2 := fmt.Sprintf("%s:%s", userID2, userID1)

	r.mx.RLock()
	defer r.mx.RUnlock()

	if friendship, ok := r.friendships[key1]; ok {
		r.logger.Info("friendship found", slog.String("user_id", userID1), slog.String("friend_id", userID2))
		return friendship, nil
	}

	if friendship, ok := r.friendships[key2]; ok {
		r.logger.Info("friendship found", slog.String("user_id", userID2), slog.String("friend_id", userID1))
		return friendship, nil
	}

	r.logger.Info("friendship not found", slog.String("user_id", userID1), slog.String("friend_id", userID2))
	return nil, nil
}

func (r *FriendshipRepository) GetAllFriendships(ctx context.Context, userID string) ([]*models.Friendship, error) {
	r.logger.Info("getting all friendships", slog.String("user_id", userID))
	r.mx.RLock()
	defer r.mx.RUnlock()

	var result []*models.Friendship
	for _, friendship := range r.friendships {
		if friendship.UserID() == userID || friendship.FriendID() == userID {
			result = append(result, friendship)
		}
	}

	r.logger.Info("all friendships found", slog.String("user_id", userID))

	return result, nil
}

func (r *FriendshipRepository) Update(ctx context.Context, friendship *models.Friendship) error {
	r.logger.Info("updating friendship",
		slog.String("user_id", friendship.UserID()), slog.String("friend_id", friendship.FriendID()))

	key := fmt.Sprintf("%s:%s", friendship.UserID(), friendship.FriendID())

	r.mx.Lock()
	defer r.mx.Unlock()

	if _, exists := r.friendships[key]; !exists {
		reverseKey := fmt.Sprintf("%s:%s", friendship.FriendID(), friendship.UserID())
		if _, exists := r.friendships[reverseKey]; !exists {
			return errors.NewNotFoundError("friendship not found", friendship.UserID(), friendship.FriendID())
		}
		key = reverseKey
	}

	r.friendships[key] = friendship

	r.logger.Info("friendship updated",
		slog.String("user_id", friendship.UserID()), slog.String("friend_id", friendship.FriendID()))
	return nil
}

func (r *FriendshipRepository) Delete(ctx context.Context, friendship *models.Friendship) error {
	r.logger.Info("deleting friendship",
		slog.String("user_id", friendship.UserID()), slog.String("friend_id", friendship.FriendID()))
	r.mx.Lock()
	defer r.mx.Unlock()

	key := fmt.Sprintf("%s:%s", friendship.UserID(), friendship.FriendID())

	if _, exists := r.friendships[key]; !exists {
		reverseKey := fmt.Sprintf("%s:%s", friendship.FriendID(), friendship.UserID())
		if _, exists := r.friendships[reverseKey]; !exists {
			return errors.NewNotFoundError("friendship not found", friendship.UserID(), friendship.FriendID())
		}
		key = reverseKey
	}

	delete(r.friendships, key)

	r.logger.Info("friendship deleted",
		slog.String("user_id", friendship.UserID()), slog.String("friend_id", friendship.FriendID()))
	return nil
}

func (r *FriendshipRepository) CheckFriendships(ctx context.Context, userIDs []string) ([]ports.UserPair, bool, error) {
	if len(userIDs) < 2 {
		return nil, true, errors.NewInvalidInputError("user IDs must contain at least 2 elements", "userIDs", userIDs)
	}

	mainUserID := userIDs[0]
	otherUserIDs := userIDs[1:]

	r.mx.RLock()
	defer r.mx.RUnlock()

	var nonFriendPairs []ports.UserPair
	for _, otherUserID := range otherUserIDs {
		key1 := fmt.Sprintf("%s:%s", mainUserID, otherUserID)
		key2 := fmt.Sprintf("%s:%s", otherUserID, mainUserID)

		friendship1, exists1 := r.friendships[key1]
		friendship2, exists2 := r.friendships[key2]

		isFriend := (exists1 && friendship1.IsAccepted()) ||
			(exists2 && friendship2.IsAccepted())

		if !isFriend {
			nonFriendPairs = append(nonFriendPairs, ports.UserPair{
				UserID1: mainUserID,
				UserID2: otherUserID,
			})
		}
	}

	allAreFriends := len(nonFriendPairs) == 0
	r.logger.Info("friendships checked", slog.String("user_id", mainUserID),
		slog.String("friend_ids", fmt.Sprintf("%v", otherUserIDs)),
		slog.Bool("all_are_friends", allAreFriends),
		slog.String("non_friend_pairs", fmt.Sprintf("%v", nonFriendPairs)),
	)
	return nonFriendPairs, allAreFriends, nil
}
