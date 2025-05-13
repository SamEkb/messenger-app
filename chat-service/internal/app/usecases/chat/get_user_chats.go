package chat

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
)

func (u *UseCase) GetUserChats(ctx context.Context, userID string) ([]*ports.ChatDto, error) {
	u.logger.Info("getting user chats", "userID", userID)

	chats, err := u.chatRepository.Get(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get user chats", "error", err)
		return nil, err
	}

	chatDtos := mapChatsToDto(ctx, chats, u)

	u.logger.Info("user chats retrieved successfully", "userID", userID)
	return chatDtos, nil
}
