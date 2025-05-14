package in_memory

import (
	"context"
	"sync"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/auth-service/pkg/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
)

var _ ports.TokenRepository = (*TokenRepository)(nil)

type TokenRepository struct {
	mx     sync.Mutex
	tokens map[models.Token]*models.AuthToken
	logger logger.Logger
}

func NewTokenRepository(logger logger.Logger) *TokenRepository {
	return &TokenRepository{
		tokens: make(map[models.Token]*models.AuthToken),
		logger: logger.With("component", "token_repository"),
	}
}

func (t *TokenRepository) Create(ctx context.Context, auth *models.AuthToken) (*models.AuthToken, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.logger.Debug("attempting to create token", "token", auth.Token(), "user_id", auth.UserID())

	token, ok := t.tokens[auth.Token()]
	if ok {
		t.logger.Warn("token already exists", "token", auth.Token())
		return nil, errors.NewAlreadyExistsError("token already exists").
			WithDetails("token", auth.Token()).
			WithDetails("user_id", auth.UserID().String())
	}

	token, err := models.NewAuthToken(auth.Token(), auth.UserID(), auth.ExpiresAt())
	if err != nil {
		t.logger.Error("failed to create token", "error", err)
		return nil, err
	}

	t.tokens[token.Token()] = token
	t.logger.Info("token created successfully", "token", token.Token(), "user_id", token.UserID(), "expires_at", token.ExpiresAt())

	return token, nil
}

func (t *TokenRepository) FindToken(ctx context.Context, token string) (models.Token, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.logger.Debug("looking for token", "token", token)

	authToken, ok := t.tokens[models.Token(token)]
	if !ok {
		t.logger.Debug("token not found", "token", token)
		return "", errors.NewNotFoundError("token not found").
			WithDetails("token", token)
	}

	t.logger.Debug("token found", "token", token, "user_id", authToken.UserID())
	return authToken.Token(), nil
}

func (t *TokenRepository) ValidateToken(ctx context.Context, token string) (bool, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.logger.Debug("validating token", "token", token)

	authToken, ok := t.tokens[models.Token(token)]
	if !ok {
		t.logger.Debug("token not found during validation", "token", token)
		return false, errors.NewNotFoundError("token not found").
			WithDetails("token", token)
	}

	if authToken.IsExpired() {
		t.logger.Debug("token is expired", "token", token, "expires_at", authToken.ExpiresAt())
		return false, errors.NewTokenError(errors.ErrTokenExpired, "token is expired").
			WithDetails("token", token).
			WithDetails("expires_at", authToken.ExpiresAt())
	}

	t.logger.Debug("token is valid", "token", token, "user_id", authToken.UserID())
	return true, nil
}

func (t *TokenRepository) DeleteToken(ctx context.Context, token models.Token) error {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.logger.Debug("attempting to delete token", "token", token)

	_, ok := t.tokens[token]
	if !ok {
		t.logger.Debug("token not found for deletion", "token", token)
		return errors.NewNotFoundError("token not found").
			WithDetails("token", token)
	}

	delete(t.tokens, token)
	t.logger.Info("token deleted successfully", "token", token)
	return nil
}
