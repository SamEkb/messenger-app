package kafka

import (
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

type UsersServiceServer struct {
	userUseCase ports.UserUseCase
	logger      logger.Logger
}

func NewUsersServiceServer(userUseCase ports.UserUseCase, logger logger.Logger) *UsersServiceServer {
	return &UsersServiceServer{
		userUseCase: userUseCase,
		logger:      logger,
	}
}
