package chat

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
)

func (u *UseCase) GetChatHistory(ctx context.Context, chatID string) ([]*ports.MessageDto, error) {
	u.logger.Info("getting chat history", "chatID", chatID)

	id, err := models.ParseChatID(chatID)
	if err != nil {
		u.logger.Error("failed to parse chat ID", "chatID", chatID, "error", err)
		return nil, err
	}
	messages, err := u.chatRepository.GetMessages(ctx, id)
	if err != nil {
		u.logger.Error("failed to get chat history", "chatID", chatID, "error", err)
		return nil, err
	}
	messageDtos := mapMessagesToDto(messages)

	u.logger.Info("chat history retrieved successfully", "chatID", chatID)
	return messageDtos, nil
}
