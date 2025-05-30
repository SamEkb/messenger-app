package in_memory

import (
	"context"
	"sync"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
)

type ChatRepository struct {
	mx       sync.Mutex
	storage  map[models.ChatID]*models.Chat
	messages map[models.ChatID][]*models.Message
	logger   logger.Logger
}

func NewChatRepository(logger logger.Logger) *ChatRepository {
	return &ChatRepository{
		storage:  make(map[models.ChatID]*models.Chat),
		messages: make(map[models.ChatID][]*models.Message),
		logger:   logger,
	}
}

func (r *ChatRepository) Create(ctx context.Context, participants []string) (*models.Chat, error) {
	r.logger.Info("creating chat", "participants", participants)

	chat, err := models.NewChat(participants)
	if err != nil {
		r.logger.Error("failed to create chat", "error", err)
		return nil, err
	}

	r.mx.Lock()
	defer r.mx.Unlock()

	r.storage[chat.ID()] = chat
	r.messages[chat.ID()] = []*models.Message{}
	r.logger.Info("chat created", "chatID", chat.ID())
	return chat, nil
}

func (r *ChatRepository) Get(ctx context.Context, userID string) ([]*models.Chat, error) {
	r.logger.Info("getting user chats", "userID", userID)

	r.mx.Lock()
	defer r.mx.Unlock()

	var result []*models.Chat
	for _, chat := range r.storage {
		for _, p := range chat.Participants() {
			if p == userID {
				result = append(result, chat)
				break
			}
		}
	}

	r.logger.Info("user chats retrieved", "chats", result)
	return result, nil
}

func (r *ChatRepository) SendMessage(ctx context.Context, chatID models.ChatID, authorID, content string) (*models.Message, error) {
	r.logger.Info("sending message", "chatID", chatID, "authorID", authorID, "content", content)

	r.mx.Lock()
	defer r.mx.Unlock()

	chat, ok := r.storage[chatID]
	if !ok {
		r.logger.Error("chat not found", "chatID", chatID)
		return nil, errors.NewNotFoundError("chat not found", "chatID", chatID)
	}
	msg, err := models.NewMessage(authorID, content)
	if err != nil {
		r.logger.Error("failed to create message", "error", err)
		return nil, err
	}

	r.messages[chatID] = append(r.messages[chatID], msg)
	if err = chat.AddMessage(msg); err != nil {
		r.logger.Error("failed to add message to chat", "error", err)
		return nil, err
	}

	r.logger.Info("message sent", "chatID", chatID, "authorID", authorID, "content", content)
	return msg, nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, chatID models.ChatID) ([]*models.Message, error) {
	r.logger.Info("getting chat history", "chatID", chatID)

	r.mx.Lock()
	defer r.mx.Unlock()

	msgs, ok := r.messages[chatID]
	if !ok {
		r.logger.Error("chat not found", "chatID", chatID)
		return nil, nil
	}

	r.logger.Info("chat history retrieved", "chatID", chatID, "messages", msgs)
	return msgs, nil
}
