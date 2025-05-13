package chat

import (
	"log/slog"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
)

var _ ports.ChatUseCase = (*UseCase)(nil)

type UseCase struct {
	chatRepository ports.ChatRepository
	userClient     ports.UserServiceClient
	friendClient   ports.FriendServiceClient
	logger         *slog.Logger
}

func NewChatUseCase(chatRepository ports.ChatRepository,
	userClient ports.UserServiceClient,
	friendClient ports.FriendServiceClient,
	logger *slog.Logger,
) *UseCase {
	return &UseCase{
		chatRepository: chatRepository,
		userClient:     userClient,
		friendClient:   friendClient,
		logger:         logger,
	}
}
