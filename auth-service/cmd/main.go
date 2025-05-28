package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/SamEkb/messenger-app/auth-service/config/env"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/adapters/in/grpc"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/adapters/out/kafka"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/repositories/auth/postgres"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	tr "github.com/SamEkb/messenger-app/pkg/platform/middleware/tracing"
	postgreslib "github.com/SamEkb/messenger-app/pkg/platform/postgres"
)

func main() {
	appCtx, cancelAppCtx := context.WithCancel(context.Background())
	defer cancelAppCtx()

	config, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(config.Debug, config.AppName)
	log.Info("starting auth service")

	tracingConfig := tr.LoadConfig()
	tracingShutdown, err := tr.Initialize(tracingConfig)
	if err != nil {
		log.Fatal("failed to initialize tracing", "error", err)
	}
	defer func() {
		if err := tracingShutdown(context.Background()); err != nil {
			log.Error("failed to shutdown tracing", "error", err)
		}
	}()

	if tracingConfig.Enabled {
		log.Info("tracing initialized", "service", tracingConfig.ServiceName, "jaeger", tracingConfig.JaegerURL)
	}

	db, err := postgreslib.NewDB(config.DB.DSN())
	if err != nil {
		log.Fatal("failed to create DB connection", "error", err)
	}
	txManager := postgreslib.NewTxManager(db)
	authRepository := postgres.NewAuthRepository(txManager, config.DB, log)
	tokenRepository := postgres.NewTokenRepository(txManager, config.DB, log)

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

	log.Info("Users service main function finished. Exiting.")
}
