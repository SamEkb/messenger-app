package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
)

type AuthUseCase interface {
	Login(ctx context.Context, dto *LoginDto) (models.Token, error)
	Register(ctx context.Context, dto *RegisterDto) (*models.User, error)
	Logout(ctx context.Context, token models.Token) error
}

type LoginDto struct {
	Email    string
	Password string
}

type RegisterDto struct {
	Username string
	Email    string
	Password string
}
