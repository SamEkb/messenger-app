package models

import (
	"errors"

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
		return nil, errors.New("id cannot be empty")
	}
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	}
	return &User{
		id:       id,
		username: username,
		email:    email,
		password: string(password),
	}, nil
}

func (u UserID) IsEmpty() bool {
	return u == UserID(uuid.Nil)
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
