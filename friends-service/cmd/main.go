package main

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	grpcserver "github.com/SamEkb/messenger-app/friends-service/internal/app/adapters/in/grpc"
	grpcclient "github.com/SamEkb/messenger-app/friends-service/internal/app/adapters/out/grpc"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/repositories/postgres"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/usecases/friendship"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	postgreslib "github.com/SamEkb/messenger-app/pkg/platform/postgres"
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

	db, err := postgreslib.NewDB(cfg.DB.DSN())
	if err != nil {
		log.Fatal("failed to create DB connection", "error", err)
	}
	txManager := postgreslib.NewTxManager(db)

	repository := postgres.NewFriendshipRepository(txManager, log)

	client := grpcclient.NewClient(cfg.Clients, log)
	usersClient, err := client.NewUsersServiceClient(ctx)
	if err != nil {
		log.Fatal("failed to create Users Service client", "error", err)
	}

	useCase := friendship.NewUseCase(repository, usersClient, txManager, log)

	server, err := grpcserver.NewServer(cfg.Server, useCase, log)
	if err != nil {
		log.Fatal("failed to create grpc server", "error", err)
	}

	if err = server.RunServers(ctx); err != nil {
		log.Fatal("failed to run grpc server", "error", err)
	}
}
