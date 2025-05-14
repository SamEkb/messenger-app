package friendship

import (
	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/pkg/platform/postgres"
)

type UseCase struct {
	friendRepository ports.FriendshipRepository
	userClient       ports.UserServiceClient
	txManager        *postgres.TxManager
	logger           logger.Logger
}

func NewUseCase(friendRepository ports.FriendshipRepository, userClient ports.UserServiceClient, txManager *postgres.TxManager, logger logger.Logger) *UseCase {
	return &UseCase{
		friendRepository: friendRepository,
		userClient:       userClient,
		txManager:        txManager,
		logger:           logger,
	}
}
