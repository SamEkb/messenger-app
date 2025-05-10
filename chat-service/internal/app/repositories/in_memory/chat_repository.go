package in_memory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/chat-service/pkg/errors"
)

type ChatRepository struct {
	mx      sync.Mutex
	storage map[models.ChatID]*models.Chat
	logger  *slog.Logger
}

func NewChatRepository(logger *slog.Logger) *ChatRepository {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	return &ChatRepository{
		storage: make(map[models.ChatID]*models.Chat),
		logger:  logger,
	}
}

func (r *ChatRepository) Create(ctx context.Context, chat *models.Chat) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Info("Creating chat", "chat", chat)

	r.storage[chat.ID()] = chat

	r.logger.Info("Chat created", "chat", chat)

	return nil
}

func (r *ChatRepository) Get(ctx context.Context, id models.ChatID) (*models.Chat, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Info("Getting chat", "id", id)

	chat, ok := r.storage[id]
	if !ok {
		r.logger.Error("Chat not found", "id", id)
		return nil, errors.NewNotFoundError("chat not found", "chatID", id)
	}

	r.logger.Info("Chat found", "chat", chat)

	return chat, nil
}

func (r *ChatRepository) GetByUserID(ctx context.Context, userID string) []*models.Chat {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Info("Getting chats by user ID", "userID", userID)

	var chats []*models.Chat
	for _, c := range r.storage {
		for _, p := range c.Participants() {
			if p == userID {
				chats = append(chats, c)
			}
		}
	}

	r.logger.Info("Chats found", "chats", chats)
	return chats
}

func (r *ChatRepository) Update(ctx context.Context, chat *models.Chat) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Info("Updating chat", "chat", chat)

	_, ok := r.storage[chat.ID()]
	if !ok {
		r.logger.Error("Chat not found", "chatID", chat.ID())
		return errors.NewNotFoundError("chat not found", "chatID", chat.ID())
	}

	r.storage[chat.ID()] = chat

	r.logger.Info("Chat updated", "chat", chat)

	return nil
}
