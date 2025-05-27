package postgres

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/pkg/platform/postgres"
)

var _ ports.TokenRepository = (*TokenRepository)(nil)

type TokenRepository struct {
	txManager *postgres.TxManager
	logger    logger.Logger
}

func NewTokenRepository(txManager *postgres.TxManager, logger logger.Logger) *TokenRepository {
	return &TokenRepository{txManager: txManager, logger: logger.With("component", "token_repository")}
}

func (r *TokenRepository) Create(ctx context.Context, token *models.AuthToken) (*models.AuthToken, error) {
	r.logger.Debug("attempting to create token", "token", token.Token(), "user_id", token.UserID())
	q := r.txManager.GetQueryEngine(ctx)
	_, err := q.ExecContext(ctx, `
		INSERT INTO tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, token.Token(), token.UserID(), token.ExpiresAt())
	if err != nil {
		r.logger.Warn("token already exists", "token", token.Token())
		return nil, errors.NewAlreadyExistsError("token already exists").WithDetails("token", token.Token())
	}
	r.logger.Debug("token created successfully", "token", token.Token(), "user_id", token.UserID())
	return token, nil
}

func (r *TokenRepository) FindToken(ctx context.Context, token string) (models.Token, error) {
	r.logger.Debug("looking for token", "token", token)
	q := r.txManager.GetQueryEngine(ctx)
	var t string
	err := q.GetContext(ctx, &t, `
		SELECT token FROM tokens WHERE token = $1
	`, token)
	if err != nil {
		r.logger.Warn("token not found", "token", token)
		return "", errors.NewNotFoundError("token not found").WithDetails("token", token)
	}
	r.logger.Debug("token found", "token", token)
	return models.Token(t), nil
}

func (r *TokenRepository) ValidateToken(ctx context.Context, token string) (bool, error) {
	r.logger.Debug("validating token", "token", token)
	q := r.txManager.GetQueryEngine(ctx)
	var expiresAt time.Time
	err := q.GetContext(ctx, &expiresAt, `
		SELECT expires_at FROM tokens WHERE token = $1
	`, token)
	if err != nil {
		r.logger.Warn("token not found during validation", "token", token)
		return false, errors.NewNotFoundError("token not found").WithDetails("token", token)
	}
	if time.Now().After(expiresAt) {
		r.logger.Debug("token is expired", "token", token, "expires_at", expiresAt)
		return false, errors.NewTokenError(errors.ErrTokenExpired, "token is expired").WithDetails("token", token).WithDetails("expires_at", expiresAt)
	}
	r.logger.Debug("token is valid", "token", token)
	return true, nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, token models.Token) error {
	r.logger.Debug("attempting to delete token", "token", token)
	q := r.txManager.GetQueryEngine(ctx)
	res, err := q.ExecContext(ctx, `
		DELETE FROM tokens WHERE token = $1
	`, token)
	if err != nil {
		r.logger.Warn("token not found for deletion", "token", token)
		return errors.NewNotFoundError("token not found").WithDetails("token", token)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		r.logger.Warn("token not found for deletion", "token", token)
		return errors.NewNotFoundError("token not found").WithDetails("token", token)
	}
	r.logger.Info("token deleted successfully", "token", token)
	return nil
}
