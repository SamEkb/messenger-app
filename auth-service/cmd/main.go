package main

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/adapters/out/kafka"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/repositories/auth/in_memory"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	authRepository := in_memory.NewAuthRepository()
	tokenRepository := in_memory.NewTokenRepository()

	userEventPublisher, err := kafka.NewUserEventsKafkaProducer()
	if err != nil {
		panic(err)
	}
	defer userEventPublisher.Close()

	usecase := auth.NewAuthUseCase(authRepository, tokenRepository, userEventPublisher, time.Hour)

	_ := usecase
}
