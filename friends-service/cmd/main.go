package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	grpcserver "github.com/SamEkb/messenger-app/friends-service/internal/app/adapters/in/grpc"
	grpcclient "github.com/SamEkb/messenger-app/friends-service/internal/app/adapters/out/grpc"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/repositories/postgres"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/usecases/friendship"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	tr "github.com/SamEkb/messenger-app/pkg/platform/middleware/tracing"
	postgreslib "github.com/SamEkb/messenger-app/pkg/platform/postgres"
	_ "github.com/lib/pq"
)

func main() {
	appCtx, cancelAppCtx := context.WithCancel(context.Background())
	defer cancelAppCtx()

	cfg, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.Debug, cfg.AppName)
	log.Info("starting friends service")

	tracingConfig := tr.LoadConfig()
	tracingShutdown, err := tr.Initialize(tracingConfig)
	if err != nil {
		log.ErrorContext(appCtx, "failed to initialize tracing", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tracingShutdown(context.Background()); err != nil {
			log.Error("failed to shutdown tracing", "error", err)
		}
	}()

	if tracingConfig.Enabled {
		log.Info("tracing initialized", "service", tracingConfig.ServiceName, "jaeger", tracingConfig.JaegerURL)
	}

	db, err := postgreslib.NewDB(cfg.DB.DSN())
	if err != nil {
		log.ErrorContext(appCtx, "failed to create DB connection", "error", err)
		os.Exit(1)
	}
	txManager := postgreslib.NewTxManager(db)

	repository := postgres.NewFriendshipRepository(txManager, cfg.DB, log)

	client := grpcclient.NewClient(cfg.Clients, log)
	usersClient, err := client.NewUsersServiceClient(appCtx)
	if err != nil {
		log.ErrorContext(appCtx, "failed to create Users Service client", "error", err)
		os.Exit(1)
	}

	useCase := friendship.NewUseCase(repository, usersClient, txManager, log)

	server, err := grpcserver.NewServer(cfg.Server, useCase, log)
	if err != nil {
		log.ErrorContext(appCtx, "failed to create grpc server", "error", err)
		os.Exit(1)
	}
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-osSignals
		log.Info("Received OS signal, initiating graceful shutdown...", "signal", sig.String())
		cancelAppCtx()
	}()

	if runErr := server.RunServers(appCtx); runErr != nil {
		if errors.Is(runErr, context.Canceled) {
			log.Info("gRPC server shutdown gracefully: context canceled.")
		} else {
			log.Error("gRPC server failed or stopped unexpectedly", "error", runErr)
		}
	} else {
		log.Info("gRPC server has shut down (RunServers returned nil).")
	}

	log.Info("Friends service main function finished. Exiting.")
}
