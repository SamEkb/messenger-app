package models

import (
	"time"

	"github.com/SamEkb/messenger-app/friends-service/pkg/errors"
	"github.com/google/uuid"
)

type FriendshipStatus string

const (
	FriendshipStatusRequested FriendshipStatus = "REQUESTED"
	FriendshipStatusAccepted  FriendshipStatus = "ACCEPTED"
	FriendshipStatusRejected  FriendshipStatus = "REJECTED"
)

type Friendship struct {
	id          uuid.UUID
	requestorID string
	recipientID string
	status      FriendshipStatus
	createdAt   time.Time
	updatedAt   time.Time
}

func NewFriendship(requestorID, recipientID string) (*Friendship, error) {
	if requestorID == "" {
		return nil, errors.NewInvalidInputError("requestorID cannot be empty")
	}
	if recipientID == "" {
		return nil, errors.NewInvalidInputError("recipientID cannot be empty")
	}

	now := time.Now()
	return &Friendship{
		id:          uuid.New(),
		requestorID: requestorID,
		recipientID: recipientID,
		status:      FriendshipStatusRequested,
		createdAt:   now,
		updatedAt:   now,
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

func (f *Friendship) RequestorID() string {
	return f.requestorID
}

func (f *Friendship) RecipientID() string {
	return f.recipientID
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
