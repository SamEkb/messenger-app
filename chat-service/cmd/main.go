package main

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	grpcserver "github.com/SamEkb/messenger-app/chat-service/internal/app/adapters/in/grpc"
	grpcclient "github.com/SamEkb/messenger-app/chat-service/internal/app/adapters/out/grpc"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/repositories/in_memory"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/usecases/chat"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(config.Debug, config.AppName)
	log.Info("starting chat service")

	chatRepository := in_memory.NewChatRepository(log)

	client := grpcclient.NewClient(config.Clients, log)

	usersClient, err := client.NewUsersServiceClient(ctx)
	if err != nil {
		log.Fatal("failed to create Users Service client", "error", err)
	}
	friendsClient, err := client.NewFriendsServiceClient(ctx)
	if err != nil {
		log.Fatal("failed to create Friends Service client", "error", err)
	}

	chatUseCase := chat.NewChatUseCase(chatRepository, usersClient, friendsClient, log)

	server, err := grpcserver.NewChatServer(chatUseCase, config.Server, log)
	if err != nil {
		log.Fatal("failed to create grpc server", "error", err)
	}

	if err = server.RunServers(ctx); err != nil {
		log.Fatal("failed to run grpc server", "error", err)
	}
}
