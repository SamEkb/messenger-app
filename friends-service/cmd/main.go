package main

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	grpcserver "github.com/SamEkb/messenger-app/friends-service/internal/app/adapters/in/grpc"
	grpcclient "github.com/SamEkb/messenger-app/friends-service/internal/app/adapters/out/grpc"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/repositories/in_memory"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/usecases/friendship"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.Debug, cfg.AppName)
	log.Info("starting friends service")

	repository := in_memory.NewFriendshipRepository(log)

	client := grpcclient.NewClient(cfg.Clients, log)
	usersClient, err := client.NewUsersServiceClient(ctx)
	if err != nil {
		log.Error("failed to create Users Service client", "error", err)
		panic(err)
	}

	useCase := friendship.NewUseCase(repository, usersClient, log)

	server, err := grpcserver.NewServer(cfg.Server, useCase, log)
	if err != nil {
		log.Error("failed to create grpc server", "error", err)
		panic(err)
	}

	if err = server.RunServers(ctx); err != nil {
		log.Error("failed to run grpc server", "error", err)
		panic(err)
	}
}
