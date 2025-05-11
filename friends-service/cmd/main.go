package main

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/repositories/in_memory"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/usecases/friendship"
	"github.com/SamEkb/messenger-app/friends-service/pkg/logger"
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

	useCase := friendship.NewUseCase(repository, log)

}
