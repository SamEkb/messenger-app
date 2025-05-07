package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/adapters/out/kafka"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/repositories/auth/in_memory"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := setupLogger(config.Debug)

	authRepository := in_memory.NewAuthRepository()
	tokenRepository := in_memory.NewTokenRepository()

	userEventPublisher, err := kafka.NewUserEventsKafkaProducer(nil, config.Kafka, logger)
	if err != nil {
		panic(err)
	}
	defer userEventPublisher.Close()

	usecase := auth.NewAuthUseCase(authRepository, tokenRepository, userEventPublisher, config.Auth.TokenTTL)

	_ := usecase
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
