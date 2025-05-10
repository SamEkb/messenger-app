package auth

import (
	"log/slog"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
)

type UseCase struct {
	authRepo           ports.AuthRepository
	tokenRepo          ports.TokenRepository
	userEventPublisher ports.UserEventsKafkaProducer
	tokenTTL           time.Duration
	logger             *slog.Logger
}

func NewAuthUseCase(
	authRepo ports.AuthRepository,
	tokenRepo ports.TokenRepository,
	userEventPublisher ports.UserEventsKafkaProducer,
	tokenTTL time.Duration,
	logger *slog.Logger,
) *UseCase {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	return &UseCase{
		authRepo:           authRepo,
		tokenRepo:          tokenRepo,
		userEventPublisher: userEventPublisher,
		tokenTTL:           tokenTTL,
		logger:             logger.With("component", "auth_usecase"),
	}
}
