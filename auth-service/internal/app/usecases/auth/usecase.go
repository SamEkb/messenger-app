package auth

import (
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
)

type UseCase struct {
	authRepo           ports.AuthRepository
	tokenRepo          ports.TokenRepository
	userEventPublisher ports.UserEventsKafkaProducer
	tokenTTL           time.Duration
}

func NewAuthUseCase(
	authRepo ports.AuthRepository,
	tokenRepo ports.TokenRepository,
	userEventPublisher ports.UserEventsKafkaProducer,
	tokenTTL time.Duration,
) *UseCase {
	return &UseCase{
		authRepo:           authRepo,
		tokenRepo:          tokenRepo,
		userEventPublisher: userEventPublisher,
		tokenTTL:           tokenTTL,
	}
}
