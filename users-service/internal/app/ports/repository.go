package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (models.UserID, error)
	Get(ctx context.Context, id models.UserID) (*models.User, error)
	GetByNickname(ctx context.Context, nickname string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}
