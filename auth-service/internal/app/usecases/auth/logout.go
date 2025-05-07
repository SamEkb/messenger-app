package auth

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
)

func (a *UseCase) Logout(ctx context.Context, token models.Token) error {
	if err := a.tokenRepo.DeleteToken(ctx, token); err != nil {
		return err
	}
	return nil
}
