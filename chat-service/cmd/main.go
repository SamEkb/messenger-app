package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	grpcserver "github.com/SamEkb/messenger-app/chat-service/internal/app/adapters/in/grpc"
	grpcclient "github.com/SamEkb/messenger-app/chat-service/internal/app/adapters/out/grpc"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/repositories/mongodb"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/usecases/chat"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	mongolib "github.com/SamEkb/messenger-app/pkg/platform/mongodb"
)

func main() {
	appCtx, cancelAppCtx := context.WithCancel(context.Background())

	config, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(config.Debug, config.AppName)
	log.Info("starting chat service")

	mongoClient, err := mongolib.NewMongoClient(appCtx, config.MongoDB.URI)
	if err != nil {
		log.Fatal("failed to connect to MongoDB", "error", err)
	}

	defer func() {
		log.Info("Disconnecting from MongoDB...")
		disconnectCtx, cancelDisconnect := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelDisconnect()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			log.Error("Failed to disconnect from MongoDB during cleanup", "error", err)
		} else {
			log.Info("MongoDB disconnected successfully.")
		}
	}()

	chatRepository := mongodb.NewChatRepository(mongoClient, config.MongoDB, log)
	txManager := mongolib.NewTxManager(mongoClient)
	client := grpcclient.NewClient(config.Clients, log)

	usersClient, err := client.NewUsersServiceClient(appCtx)
	if err != nil {
		log.Fatal("failed to create Users Service client", "error", err)
	}
	friendsClient, err := client.NewFriendsServiceClient(appCtx)
	if err != nil {
		log.Fatal("failed to create Friends Service client", "error", err)
	}

	defer func() {
		log.Info("Closing gRPC users client...")
		if err := usersClient.Close(); err != nil {
			log.Error("Failed to close users client", "error", err)
		} else {
			log.Info("Users client closed successfully")
		}
	}()

	defer func() {
		log.Info("Closing gRPC friends client...")
		if err := friendsClient.Close(); err != nil {
			log.Error("Failed to close friends client", "error", err)
		} else {
			log.Info("Friends client closed successfully")
		}
	}()

	chatUseCase := chat.NewChatUseCase(chatRepository, usersClient, friendsClient, txManager, log)

	server, err := grpcserver.NewChatServer(chatUseCase, config.Server, log)
	if err != nil {
		log.Fatal("failed to create grpc server", "error", err)
	}

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-osSignals
		log.Info("Received OS signal, initiating graceful shutdown...", "signal", sig.String())
		cancelAppCtx()
	}()

	log.Info("gRPC server starting...", "address", config.Server.GRPCPort)

	if runErr := server.RunServers(appCtx); runErr != nil {
		if errors.Is(runErr, context.Canceled) {
			log.Info("gRPC server shutdown gracefully: context canceled.")
		} else if runErr.Error() == "grpc: the server has been stopped" {
			log.Info("gRPC server shutdown gracefully: server stopped.")
		} else {
			log.Error("gRPC server failed or stopped unexpectedly", "error", runErr)
		}
	} else {
		log.Info("gRPC server has shut down (RunServers returned nil).")
	}

	log.Info("Chat service main function finished. Exiting.")
}
