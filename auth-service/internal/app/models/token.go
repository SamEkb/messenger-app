package models

import (
	"errors"
	"time"
)

type Token string

type AuthToken struct {
	token     Token
	userID    UserID
	expiresAt time.Time
}

func NewAuthToken(token Token, userID UserID, expiresAt time.Time) (*AuthToken, error) {
	if token.IsEmpty() {
		return nil, errors.New("token cannot be empty")
	}

	if expiresAt.IsZero() {
		return nil, errors.New("expiresAt cannot be zero")
	}

	if userID.IsEmpty() {
		return nil, errors.New("userID cannot be empty")
	}

	return &AuthToken{
		token:     token,
		userID:    userID,
		expiresAt: expiresAt,
	}, nil
}

func (a *AuthToken) Token() Token {
	return a.token
}

func (a *AuthToken) UserID() UserID {
	return a.userID
}

func (a *AuthToken) ExpiresAt() time.Time {
	return a.expiresAt
}

func (t Token) IsEmpty() bool {
	return t == ""
}

func (a *AuthToken) IsExpired() bool {
	return time.Now().After(a.expiresAt)
}
