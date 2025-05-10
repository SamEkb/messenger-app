package kafka

import (
	"log/slog"

	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

type UsersServiceServer struct {
	userUseCase ports.UserUseCase
	logger      *slog.Logger
}

func NewUsersServiceServer(userUseCase ports.UserUseCase, logger *slog.Logger) *UsersServiceServer {
	return &UsersServiceServer{
		userUseCase: userUseCase,
		logger:      logger,
	}
}
