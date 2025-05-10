package models

import (
	"time"

	"github.com/SamEkb/messenger-app/friends-service/pkg/errors"
	"github.com/google/uuid"
)

// FriendshipStatus represents the status of friendship between two users
type FriendshipStatus string

const (
	FriendshipStatusRequested FriendshipStatus = "REQUESTED"
	FriendshipStatusAccepted  FriendshipStatus = "ACCEPTED"
	FriendshipStatusRejected  FriendshipStatus = "REJECTED"
)

type Friendship struct {
	id        uuid.UUID
	userID    string
	friendID  string
	status    FriendshipStatus
	createdAt time.Time
	updatedAt time.Time
}

func NewFriendship(userID, friendID string) (*Friendship, error) {
	if userID == "" {
		return nil, errors.NewInvalidInputError("userID cannot be empty")
	}
	if friendID == "" {
		return nil, errors.NewInvalidInputError("friendID cannot be empty")
	}

	now := time.Now()
	return &Friendship{
		id:        uuid.New(),
		userID:    userID,
		friendID:  friendID,
		status:    FriendshipStatusRequested,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (f *Friendship) Accept() {
	f.status = FriendshipStatusAccepted
	f.updatedAt = time.Now()
}

func (f *Friendship) Reject() {
	f.status = FriendshipStatusRejected
	f.updatedAt = time.Now()
}

func (f *Friendship) IsAccepted() bool {
	return f.status == FriendshipStatusAccepted
}

func (f *Friendship) IsRequested() bool {
	return f.status == FriendshipStatusRequested
}

func (f *Friendship) IsRejected() bool {
	return f.status == FriendshipStatusRejected
}

func (f *Friendship) ID() uuid.UUID {
	return f.id
}

func (f *Friendship) UserID() string {
	return f.userID
}

func (f *Friendship) FriendID() string {
	return f.friendID
}

func (f *Friendship) Status() FriendshipStatus {
	return f.status
}

func (f *Friendship) CreatedAt() time.Time {
	return f.createdAt
}

func (f *Friendship) UpdatedAt() time.Time {
	return f.updatedAt
}
