package ports

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/friends-service/pkg/errors"
)

type FriendshipUseCase interface {
	GetFriends(ctx context.Context, userID string) ([]*FriendshipDto, error)
	SendFriendRequest(ctx context.Context, requestorID, recipientID string) error
	AcceptFriendRequest(ctx context.Context, recipientID, requestorID string) error
	RejectFriendRequest(ctx context.Context, recipientID, requestorID string) error
	DeleteFriend(ctx context.Context, userID string, friendID string) error
}

type FriendshipDto struct {
	id               string
	requestorID      string
	recipientID      string
	friendsNickName  string
	friendsAvatarURL string
	status           string
	createdAt        time.Time
	updatedAt        time.Time
}

func NewFriendshipDto(id, requestorID, recipientID, friendsNickName, friendsAvatarURL, status string,
	createdAt, updatedAt time.Time) (*FriendshipDto, error) {
	if id == "" {
		return nil, errors.NewInvalidInputError("id cannot be empty")
	}
	if requestorID == "" {
		return nil, errors.NewInvalidInputError("requestorID cannot be empty")
	}
	if recipientID == "" {
		return nil, errors.NewInvalidInputError("recipientID cannot be empty")
	}
	if friendsNickName == "" {
		return nil, errors.NewInvalidInputError("friendsNickName cannot be empty")
	}
	if status == "" {
		return nil, errors.NewInvalidInputError("status cannot be empty")
	}
	if createdAt.IsZero() {
		return nil, errors.NewInvalidInputError("createdAt cannot be zero")
	}
	if updatedAt.IsZero() {
		return nil, errors.NewInvalidInputError("updatedAt cannot be zero")
	}

	return &FriendshipDto{
		id:               id,
		requestorID:      requestorID,
		recipientID:      recipientID,
		friendsNickName:  friendsNickName,
		friendsAvatarURL: friendsAvatarURL,
		status:           status,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}, nil
}

func (f *FriendshipDto) ID() string {
	return f.id
}

func (f *FriendshipDto) RequestorID() string {
	return f.requestorID
}

func (f *FriendshipDto) RecipientID() string {
	return f.recipientID
}

func (f *FriendshipDto) FriendsNickName() string {
	return f.friendsNickName
}

func (f *FriendshipDto) FriendsAvatarURL() string {
	return f.friendsAvatarURL
}

func (f *FriendshipDto) Status() string {
	return f.status
}

func (f *FriendshipDto) CreatedAt() time.Time {
	return f.createdAt
}

func (f *FriendshipDto) UpdatedAt() time.Time {
	return f.updatedAt
}
