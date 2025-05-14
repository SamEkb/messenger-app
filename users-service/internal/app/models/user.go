package models

import (
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/google/uuid"
)

type UserID uuid.UUID

type User struct {
	id          UserID
	email       string
	nickname    string
	description string
	avatarUrl   string
}

func NewUser(id UserID, email string, nickname string, description string, avatarUrl string) (*User, error) {
	if id.IsEmpty() {
		return nil, errors.NewInvalidInputError("id cannot be empty").
			WithDetails("id", "nickname")
	}
	if nickname == "" {
		return nil, errors.NewInvalidInputError("nickname cannot be empty").
			WithDetails("field", "nickname")
	}

	if email == "" {
		return nil, errors.NewInvalidInputError("email cannot be empty").
			WithDetails("field", "email")
	}

	return &User{
		id:          id,
		nickname:    nickname,
		email:       email,
		description: description,
		avatarUrl:   avatarUrl,
	}, nil
}

func (u UserID) IsEmpty() bool {
	return u == UserID(uuid.Nil)
}

func (u UserID) String() string {
	return uuid.UUID(u).String()
}

func ParseUserID(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, errors.NewInvalidInputError("invalid UUID format: %s", s).
			WithDetails("value", s)
	}
	return UserID(id), nil
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Nickname() string {
	return u.nickname
}

func (u *User) Description() string {
	return u.description
}

func (u *User) AvatarURL() string {
	return u.avatarUrl
}
