package main

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	grpcserver "github.com/SamEkb/messenger-app/chat-service/internal/app/adapters/in/grpc"
	grpcclient "github.com/SamEkb/messenger-app/chat-service/internal/app/adapters/out/grpc"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/repositories/mongodb"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/usecases/chat"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	mongolib "github.com/SamEkb/messenger-app/pkg/platform/mongodb"
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

	mongoClient, err := mongolib.NewMongoClient(ctx, config.MongoDB.URI)
	if err != nil {
		log.Fatal("failed to connect to MongoDB", "error", err)
	}
	defer mongoClient.Disconnect(ctx)

	chatRepository := mongodb.NewChatRepository(mongoClient, config.MongoDB, log)

	txManager := mongolib.NewTxManager(mongoClient)

	client := grpcclient.NewClient(config.Clients, log)

	usersClient, err := client.NewUsersServiceClient(ctx)
	if err != nil {
		log.Fatal("failed to create Users Service client", "error", err)
	}
	friendsClient, err := client.NewFriendsServiceClient(ctx)
	if err != nil {
		log.Fatal("failed to create Friends Service client", "error", err)
	}

	chatUseCase := chat.NewChatUseCase(chatRepository, usersClient, friendsClient, txManager, log)

	server, err := grpcserver.NewChatServer(chatUseCase, config.Server, log)
	if err != nil {
		log.Fatal("failed to create grpc server", "error", err)
	}

	if err = server.RunServers(ctx); err != nil {
		log.Fatal("failed to run grpc server", "error", err)
	}
}
