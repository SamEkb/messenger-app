package main

import (
	"context"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	postgreslib "github.com/SamEkb/messenger-app/pkg/platform/postgres"
	"github.com/SamEkb/messenger-app/users-service/config/env"
	"github.com/SamEkb/messenger-app/users-service/internal/app/adapters/in/grpc"
	"github.com/SamEkb/messenger-app/users-service/internal/app/adapters/in/kafka"
	"github.com/SamEkb/messenger-app/users-service/internal/app/repositories/user/postgres"
	"github.com/SamEkb/messenger-app/users-service/internal/app/usecases/user"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(config.Debug, config.AppName)
	log.Info("starting users service")

	db, err := postgreslib.NewDB(config.DB.DSN())
	if err != nil {
		log.Fatal("failed to create DB connection", "error", err)
	}
	txManager := postgreslib.NewTxManager(db)

	usersRepo := postgres.NewUserRepository(txManager, log)
	userUseCase := user.NewUseCase(usersRepo, txManager, log)

	kafkaServer := kafka.NewUsersServiceServer(userUseCase, log)

	consumer, err := kafka.NewConsumerWithConfig(kafkaServer, config.Kafka)
	if err != nil {
		log.Fatal("failed to create Kafka consumer", "error", err)
	}

	if err := consumer.Start(ctx); err != nil {
		log.Fatal("failed to start Kafka consumer", "error", err)
	}
	defer consumer.Close()

	grpcServer, err := grpc.NewServer(config.Server, userUseCase, log)
	if err != nil {
		log.Fatal("failed to create grpc server", "error", err)
	}

	if err = grpcServer.RunServers(ctx); err != nil {
		log.Fatal("failed to run grpc server", "error", err)
	}
}
