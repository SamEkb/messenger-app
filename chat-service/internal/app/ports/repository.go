package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) error
	Get(ctx context.Context, id models.ChatID) (*models.Chat, error)
	GetByUserID(ctx context.Context, userID string) []*models.Chat
	Update(ctx context.Context, chat *models.Chat) error
}
