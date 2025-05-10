package user

import (
	"log/slog"

	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

var _ ports.UserUseCase = (*UseCase)(nil)

type UseCase struct {
	userRepository ports.UserRepository
	logger         *slog.Logger
}

func NewUseCase(userRepository ports.UserRepository, logger *slog.Logger) *UseCase {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	return &UseCase{
		userRepository: userRepository,
		logger:         logger.With("component", "user_usecase"),
	}
}
