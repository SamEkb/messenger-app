package main

import (
	"context"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/users-service/config/env"
	"github.com/SamEkb/messenger-app/users-service/internal/app/adapters/in/grpc"
	"github.com/SamEkb/messenger-app/users-service/internal/app/adapters/in/kafka"
	"github.com/SamEkb/messenger-app/users-service/internal/app/repositories/user/in_memory"
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

	usersRepo := in_memory.NewUserRepository(log)
	userUseCase := user.NewUseCase(usersRepo, log)

	kafkaServer := kafka.NewUsersServiceServer(userUseCase, log)

	consumer, err := kafka.NewConsumerWithConfig(kafkaServer, config.Kafka)
	if err != nil {
		log.Error("failed to create Kafka consumer", "error", err)
		panic(err)
	}

	if err := consumer.Start(ctx); err != nil {
		log.Error("failed to start Kafka consumer", "error", err)
		panic(err)
	}
	defer consumer.Close()

	grpcServer, err := grpc.NewServer(config.Server, userUseCase, log)
	if err != nil {
		log.Error("failed to create grpc server", "error", err)
		panic(err)
	}

	if err = grpcServer.RunServers(ctx); err != nil {
		log.Error("failed to run grpc server", "error", err)
		panic(err)
	}
}
