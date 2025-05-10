package ports

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
)

type ChatRepository interface {
	Create(ctx context.Context, participants []string) (*models.Chat, error)
	Get(ctx context.Context, userID string) ([]*models.Chat, error)
	SendMessage(ctx context.Context, chatID models.ChatID, authorID, content string) (*models.Message, error)
	GetMessages(ctx context.Context, chatID models.ChatID) ([]*models.Message, error)
}
