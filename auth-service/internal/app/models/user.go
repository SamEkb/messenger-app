package models

import (
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/google/uuid"
)

type UserID uuid.UUID

type User struct {
	id       UserID
	username string
	email    string
	password string
}

func NewUser(id UserID, username string, email string, password []byte) (*User, error) {
	if id.IsEmpty() {
		return nil, errors.NewInvalidInputError("id cannot be empty")
	}
	if username == "" {
		return nil, errors.NewInvalidInputError("username cannot be empty")
	}
	if email == "" {
		return nil, errors.NewInvalidInputError("email cannot be empty").
			WithDetails("field", "email")
	}
	if len(password) == 0 {
		return nil, errors.NewInvalidInputError("password cannot be empty").
			WithDetails("field", "password")
	}
	return &User{
		id:       id,
		username: username,
		email:    email,
		password: string(password),
	}, nil
}

func UserIDFromString(id string) (UserID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, errors.NewInvalidInputError("invalid user ID")
	}
	return UserID(parsed), nil
}

func (u UserID) IsEmpty() bool {
	return u == UserID(uuid.Nil)
}

func (u UserID) String() string {
	return uuid.UUID(u).String()
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}
