package friendship

import (
	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
)

type UseCase struct {
	friendRepository ports.FriendshipRepository
	userClient       ports.UserServiceClient
	logger           logger.Logger
}

func NewUseCase(friendRepository ports.FriendshipRepository, userClient ports.UserServiceClient, logger logger.Logger) *UseCase {
	return &UseCase{
		friendRepository: friendRepository,
		userClient:       userClient,
		logger:           logger,
	}
}
