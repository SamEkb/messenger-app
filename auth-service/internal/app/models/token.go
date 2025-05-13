package models

import (
	"time"

	"github.com/SamEkb/messenger-app/auth-service/pkg/errors"
)

type Token string

type AuthToken struct {
	token     Token
	userID    UserID
	expiresAt time.Time
}

func NewAuthToken(token Token, userID UserID, expiresAt time.Time) (*AuthToken, error) {
	if token.IsEmpty() {
		return nil, errors.NewInvalidInputError("token cannot be empty").
			WithDetails("field", "token")
	}

	if expiresAt.IsZero() {
		return nil, errors.NewInvalidInputError("expiresAt cannot be zero").
			WithDetails("field", "expiresAt")
	}

	if userID.IsEmpty() {
		return nil, errors.NewInvalidInputError("userID cannot be empty").
			WithDetails("field", "userID")
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
