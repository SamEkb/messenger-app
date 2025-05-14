package chat

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (u *UseCase) SendMessage(ctx context.Context, chatID string, authorID, content string) (*ports.MessageDto, error) {
	u.logger.Info("sending message", "chatID", chatID, "authorID", authorID)

	if chatID == "" {
		err := errors.NewInvalidInputError("chat ID is required")
		u.logger.Error("invalid input", "error", err)
		return nil, err
	}
	if authorID == "" {
		err := errors.NewInvalidInputError("author ID is required")
		u.logger.Error("invalid input", "error", err)
		return nil, err
	}
	if content == "" {
		err := errors.NewInvalidInputError("message content is required")
		u.logger.Error("invalid input", "error", err)
		return nil, err
	}

	_, err := u.userClient.GetUserProfile(authorID)
	if err != nil {
		u.logger.Error("failed to get user profile", "authorID", authorID, "error", err)
		return nil, err
	}

	id, err := models.ParseChatID(chatID)
	if err != nil {
		u.logger.Error("failed to parse chat ID", "chatID", chatID, "error", err)
		return nil, err
	}

	var msg *models.Message
	err = u.txManager.RunTx(ctx, func(sessionCtx mongo.SessionContext) error {
		var err error
		msg, err = u.chatRepository.SendMessage(sessionCtx, id, authorID, content)
		if err != nil {
			u.logger.Error("failed to send message", "chatID", chatID, "authorID", authorID, "error", err)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	dto := ports.NewMessageDto(msg.ID().String(), msg.AuthorID(), msg.Content(), msg.Timestamp())

	u.logger.Info("message sent successfully", "chatID", chatID, "authorID", authorID)
	return dto, nil
}
