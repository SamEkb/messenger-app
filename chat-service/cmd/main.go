package main

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/repositories/in_memory"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/usecases/chat"
	"github.com/SamEkb/messenger-app/chat-service/pkg/logger"
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

	chatUseCase := chat.NewChatUseCase()
}
