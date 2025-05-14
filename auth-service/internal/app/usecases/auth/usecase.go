package auth

import (
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
)

type UseCase struct {
	authRepo           ports.AuthRepository
	tokenRepo          ports.TokenRepository
	userEventPublisher ports.UserEventsKafkaProducer
	tokenTTL           time.Duration
	logger             logger.Logger
}

func NewAuthUseCase(
	authRepo ports.AuthRepository,
	tokenRepo ports.TokenRepository,
	userEventPublisher ports.UserEventsKafkaProducer,
	tokenTTL time.Duration,
	logger logger.Logger,
) *UseCase {
	return &UseCase{
		authRepo:           authRepo,
		tokenRepo:          tokenRepo,
		userEventPublisher: userEventPublisher,
		tokenTTL:           tokenTTL,
		logger:             logger.With("component", "auth_usecase"),
	}
}
