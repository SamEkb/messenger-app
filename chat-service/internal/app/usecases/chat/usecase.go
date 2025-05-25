package chat

import (
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/mongodb"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
)

var _ ports.ChatUseCase = (*UseCase)(nil)

type UseCase struct {
	chatRepository ports.ChatRepository
	userClient     ports.UserServiceClient
	friendClient   ports.FriendServiceClient
	txManager      *mongodb.TxManager
	logger         logger.Logger
}

func NewChatUseCase(chatRepository ports.ChatRepository,
	userClient ports.UserServiceClient,
	friendClient ports.FriendServiceClient,
	txManager *mongodb.TxManager,
	logger logger.Logger,
) *UseCase {
	return &UseCase{
		chatRepository: chatRepository,
		userClient:     userClient,
		friendClient:   friendClient,
		txManager:      txManager,
		logger:         logger,
	}
}
