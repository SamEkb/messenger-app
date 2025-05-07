package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
)

type AuthRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindUserByID(ctx context.Context, userID models.UserID) (*models.User, error)
	FindUserByEmail(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type TokenRepository interface {
	Create(ctx context.Context, token *models.AuthToken) (*models.AuthToken, error)
	FindToken(ctx context.Context, token string) (models.Token, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
	DeleteToken(ctx context.Context, token models.Token) error
}
