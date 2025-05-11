package friendship

import (
	"log/slog"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
)

type UseCase struct {
	friendRepository ports.FriendshipRepository
	userClient       ports.UserServiceClient
	logger           *slog.Logger
}

func NewUseCase(friendRepository ports.FriendshipRepository, userClient ports.UserServiceClient, logger *slog.Logger) *UseCase {
	return &UseCase{
		friendRepository: friendRepository,
		userClient:       userClient,
		logger:           logger,
	}
}
