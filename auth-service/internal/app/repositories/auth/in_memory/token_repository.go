package in_memory

import (
	"context"
	"errors"
	"sync"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
)

var _ ports.TokenRepository = (*TokenRepository)(nil)

type TokenRepository struct {
	mx     sync.Mutex
	tokens map[models.Token]*models.AuthToken
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{
		tokens: make(map[models.Token]*models.AuthToken),
	}
}

func (t *TokenRepository) Create(ctx context.Context, auth *models.AuthToken) (*models.AuthToken, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	token, ok := t.tokens[auth.Token()]
	if ok {
		return nil, errors.New("token already exists")
	}

	token, err := models.NewAuthToken(auth.Token(), auth.UserID(), auth.ExpiresAt())
	if err != nil {
		return nil, err
	}

	t.tokens[token.Token()] = token

	return token, nil
}

func (t *TokenRepository) FindToken(ctx context.Context, token string) (models.Token, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	authToken, ok := t.tokens[models.Token(token)]
	if !ok {
		return "", errors.New("token not found")
	}

	return authToken.Token(), nil
}

func (t *TokenRepository) ValidateToken(ctx context.Context, token string) (bool, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	authToken, ok := t.tokens[models.Token(token)]
	if !ok {
		return false, errors.New("token not found")
	}

	if authToken.IsExpired() {
		return false, errors.New("token is expired")
	}

	return true, nil
}

func (t *TokenRepository) DeleteToken(ctx context.Context, token models.Token) error {
	t.mx.Lock()
	defer t.mx.Unlock()
	_, ok := t.tokens[token]
	if !ok {
		return errors.New("token not found")
	}

	delete(t.tokens, token)
	return nil
}
