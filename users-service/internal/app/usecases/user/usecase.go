package user

import (
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

var _ ports.UserUseCase = (*UseCase)(nil)

type UseCase struct {
	userRepository ports.UserRepository
	logger         logger.Logger
}

func NewUseCase(userRepository ports.UserRepository, logger logger.Logger) *UseCase {
	return &UseCase{
		userRepository: userRepository,
		logger:         logger.With("component", "user_usecase"),
	}
}
