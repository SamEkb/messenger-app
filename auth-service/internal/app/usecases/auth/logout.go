package auth

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/pkg/errors"
)

func (a *UseCase) Logout(ctx context.Context, token models.Token) error {
	a.logger.Debug("logout attempt", "token", token)

	valid, err := a.tokenRepo.ValidateToken(ctx, string(token))
	if err != nil {
		a.logger.Warn("token validation failed", "error", err)
		return err
	}

	if !valid {
		a.logger.Warn("invalid token", "token", token)
		return errors.NewTokenError(errors.ErrInvalidToken, "token is invalid").
			WithDetails("token", token)
	}

	err = a.tokenRepo.DeleteToken(ctx, token)
	if err != nil {
		a.logger.Error("failed to delete token", "error", err)
		return err
	}

	a.logger.Info("logout successful", "token", token)
	return nil
}
