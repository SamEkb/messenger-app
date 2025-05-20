package main

import (
	"context"

	_ "github.com/lib/pq"

	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/adapters/in/grpc"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/adapters/out/kafka"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/repositories/auth/postgres"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	postgreslib "github.com/SamEkb/messenger-app/pkg/platform/postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(config.Debug, config.AppName)
	log.Info("starting auth service")

	db, err := postgreslib.NewDB(config.DB.DSN())
	if err != nil {
		log.Fatal("failed to create DB connection", "error", err)
	}
	txManager := postgreslib.NewTxManager(db)
	authRepository := postgres.NewAuthRepository(txManager, log)
	tokenRepository := postgres.NewTokenRepository(txManager, log)

	userEventPublisher, err := kafka.NewUserEventsKafkaProducer(config.Kafka, log)
	if err != nil {
		log.Fatal("failed to create Kafka producer", "error", err)
	}
	defer userEventPublisher.Close()

	usecase := auth.NewAuthUseCase(
		txManager,
		authRepository,
		tokenRepository,
		userEventPublisher,
		config.Auth.TokenTTL,
		log,
	)

	server, err := grpc.NewServer(config.Server, usecase, log)
	if err != nil {
		log.Fatal("failed to create grpc server", "error", err)
	}

	if err = server.RunServers(ctx); err != nil {
		log.Fatal("failed to run grpc server", "error", err)
	}
}
