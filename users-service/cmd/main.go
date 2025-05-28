package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	tr "github.com/SamEkb/messenger-app/pkg/platform/middleware/tracing"
	postgreslib "github.com/SamEkb/messenger-app/pkg/platform/postgres"
	"github.com/SamEkb/messenger-app/users-service/config/env"
	"github.com/SamEkb/messenger-app/users-service/internal/app/adapters/in/grpc"
	"github.com/SamEkb/messenger-app/users-service/internal/app/adapters/in/kafka"
	"github.com/SamEkb/messenger-app/users-service/internal/app/repositories/user/postgres"
	"github.com/SamEkb/messenger-app/users-service/internal/app/usecases/user"
)

func main() {
	appCtx, cancelAppCtx := context.WithCancel(context.Background())

	config, err := env.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(config.Debug, config.AppName)
	log.Info("starting users service")

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

	usersRepo := postgres.NewUserRepository(txManager, log, config.DB)
	userUseCase := user.NewUseCase(usersRepo, txManager, log)

	kafkaServer := kafka.NewUsersServiceServer(userUseCase, log)

	consumer, err := kafka.NewConsumerWithConfig(kafkaServer, config.Kafka)
	if err != nil {
		log.Fatal("failed to create Kafka consumer", "error", err)
	}

	if err := consumer.Start(appCtx); err != nil {
		log.Fatal("failed to start Kafka consumer", "error", err)
	}
	defer consumer.Close()

	grpcServer, err := grpc.NewServer(config.Server, userUseCase, log)
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

	if runErr := grpcServer.RunServers(appCtx); runErr != nil {
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
