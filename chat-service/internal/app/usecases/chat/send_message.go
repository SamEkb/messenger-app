package chat

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/models"
	"github.com/SamEkb/messenger-app/chat-service/pkg/errors"
)

func (u *UseCase) SendMessage(ctx context.Context, chatID string, authorID, content string) error {
	u.logger.Info("sending message", "chatID", chatID, "authorID", authorID)

	if chatID == "" {
		err := errors.NewInvalidInputError("chat ID is required")
		u.logger.Error("invalid input", "error", err)
		return err
	}
	if authorID == "" {
		err := errors.NewInvalidInputError("author ID is required")
		u.logger.Error("invalid input", "error", err)
		return err
	}
	if content == "" {
		err := errors.NewInvalidInputError("message content is required")
		u.logger.Error("invalid input", "error", err)
		return err
	}

	_, err := u.userClient.GetUserProfile(authorID)
	if err != nil {
		u.logger.Error("failed to get user profile", "authorID", authorID, "error", err)
		return err
	}

	id, err := models.ParseChatID(chatID)
	if err != nil {
		u.logger.Error("failed to parse chat ID", "chatID", chatID, "error", err)
		return err
	}

	if err = u.chatRepository.SendMessage(ctx, id, authorID, content); err != nil {
		u.logger.Error("failed to send message", "chatID", chatID, "authorID", authorID, "error", err)
		return err
	}

	u.logger.Info("message sent successfully", "chatID", chatID, "authorID", authorID)
	return nil
}
