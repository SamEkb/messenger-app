package chat

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
)

func mapChatsToDto(ctx context.Context, chats []*models.Chat, u *UseCase) []*ports.ChatDto {
	chatDtos := make([]*ports.ChatDto, 0, len(chats))

	for _, chat := range chats {
		messages, err := u.chatRepository.GetMessages(ctx, chat.ID())
		if err != nil {
			u.logger.Error("failed to get messages",
				"chatID", chat.ID().String(),
				"error", err)

			messages = make([]*models.Message, 0)
		}

		messageDtos := mapMessagesToDto(messages)

		chatDto := ports.NewChatDto(
			chat.ID().String(),
			chat.Participants(),
			messageDtos,
			chat.CreatedAt(),
			chat.UpdatedAt(),
		)

		chatDtos = append(chatDtos, chatDto)
	}

	return chatDtos
}

func mapMessagesToDto(messages []*models.Message) []*ports.MessageDto {
	messageDtos := make([]*ports.MessageDto, 0, len(messages))
	for _, message := range messages {
		messageDto := ports.NewMessageDto(
			message.ID().String(),
			message.AuthorID(),
			message.Content(),
			message.Timestamp(),
		)

		messageDtos = append(messageDtos, messageDto)
	}

	return messageDtos
}
