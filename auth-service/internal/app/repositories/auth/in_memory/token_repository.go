package in_memory

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
)

var _ ports.TokenRepository = (*TokenRepository)(nil)

type TokenRepository struct {
	mx     sync.Mutex
	tokens map[models.Token]*models.AuthToken
	logger *slog.Logger
}

func NewTokenRepository(logger *slog.Logger) *TokenRepository {
	// Если логгер не передан, создаем дефолтный noop логгер
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

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
		return nil, errors.New("token already exists")
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
		return "", errors.New("token not found")
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
		return false, errors.New("token not found")
	}

	if authToken.IsExpired() {
		t.logger.Debug("token is expired", "token", token, "expires_at", authToken.ExpiresAt())
		return false, errors.New("token is expired")
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
		return errors.New("token not found")
	}

	delete(t.tokens, token)
	t.logger.Info("token deleted successfully", "token", token)
	return nil
}
