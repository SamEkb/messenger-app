package user

import (
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/pkg/platform/postgres"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

var _ ports.UserUseCase = (*UseCase)(nil)

type UseCase struct {
	userRepository ports.UserRepository
	txManager      *postgres.TxManager
	logger         logger.Logger
}

func NewUseCase(userRepository ports.UserRepository, txManager *postgres.TxManager, logger logger.Logger) *UseCase {
	return &UseCase{
		userRepository: userRepository,
		txManager:      txManager,
		logger:         logger.With("component", "user_usecase"),
	}
}
