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
	friendships map[string][]*models.Friendship // userID -> список дружеских отношений
	mx          sync.RWMutex
	logger      *slog.Logger
}

func NewFriendshipRepository(logger *slog.Logger) *FriendshipRepository {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	return &FriendshipRepository{
		friendships: make(map[string][]*models.Friendship),
		logger:      logger.With("component", "friendship_repository"),
	}
}

func (r *FriendshipRepository) GetFriends(ctx context.Context, userID string) ([]*models.Friendship, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	friendships, exists := r.friendships[userID]
	if !exists {
		return []*models.Friendship{}, nil
	}

	acceptedFriendships := make([]*models.Friendship, 0)
	for _, friendship := range friendships {
		if friendship.IsAccepted() {
			acceptedFriendships = append(acceptedFriendships, friendship)
		}
	}

	return acceptedFriendships, nil
}

func (r *FriendshipRepository) SendFriendRequest(ctx context.Context, requestorID, recipientID string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if r.friendshipExists(requestorID, recipientID) {
		return errors.NewInvalidInputError("friendship request already exists")
	}

	friendship, err := models.NewFriendship(requestorID, recipientID)
	if err != nil {
		return err
	}

	r.addFriendshipToUser(requestorID, friendship)
	r.addFriendshipToUser(recipientID, friendship)

	return nil
}

func (r *FriendshipRepository) AcceptFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	friendship, err := r.findFriendshipRequest(recipientID, requestorID)
	if err != nil {
		return err
	}

	if !friendship.IsRequested() {
		return errors.NewInvalidInputError(fmt.Sprintf("friendship is not in REQUESTED state, current state: %s", friendship.Status()))
	}

	friendship.Accept()
	return nil
}

func (r *FriendshipRepository) RejectFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	friendship, err := r.findFriendshipRequest(recipientID, requestorID)
	if err != nil {
		return err
	}

	if !friendship.IsRequested() {
		return errors.NewInvalidInputError(fmt.Sprintf("friendship is not in REQUESTED state, current state: %s", friendship.Status()))
	}

	friendship.Reject()
	return nil
}

func (r *FriendshipRepository) Delete(ctx context.Context, userID string, friendID string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.removeFriendship(userID, friendID)
	r.removeFriendship(friendID, userID)

	return nil
}

func (r *FriendshipRepository) friendshipExists(user1ID, user2ID string) bool {
	friendships, exists := r.friendships[user1ID]
	if !exists {
		return false
	}

	for _, friendship := range friendships {
		if (friendship.RequestorID() == user1ID && friendship.RecipientID() == user2ID) ||
			(friendship.RequestorID() == user2ID && friendship.RecipientID() == user1ID) {
			return true
		}
	}

	return false
}

func (r *FriendshipRepository) addFriendshipToUser(userID string, friendship *models.Friendship) {
	if _, exists := r.friendships[userID]; !exists {
		r.friendships[userID] = make([]*models.Friendship, 0)
	}
	r.friendships[userID] = append(r.friendships[userID], friendship)
}

func (r *FriendshipRepository) findFriendshipRequest(recipientID, requestorID string) (*models.Friendship, error) {
	friendships, exists := r.friendships[recipientID]
	if !exists {
		return nil, errors.NewNotFoundError("friendship request not found")
	}

	for _, friendship := range friendships {
		if friendship.RequestorID() == requestorID && friendship.RecipientID() == recipientID {
			return friendship, nil
		}
	}

	return nil, errors.NewNotFoundError("friendship request not found")
}

func (r *FriendshipRepository) removeFriendship(userID, otherUserID string) {
	friendships, exists := r.friendships[userID]
	if !exists {
		return
	}

	updatedFriendships := make([]*models.Friendship, 0)
	for _, friendship := range friendships {
		if !((friendship.RequestorID() == userID && friendship.RecipientID() == otherUserID) ||
			(friendship.RequestorID() == otherUserID && friendship.RecipientID() == userID)) {
			updatedFriendships = append(updatedFriendships, friendship)
		}
	}

	r.friendships[userID] = updatedFriendships
}
