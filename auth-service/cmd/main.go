package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SamEkb/messenger-app/auth-service/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("received shutdown signal")
		cancel()
	}()

	serv, err := server.NewServer()
	if err != nil {
		log.Fatalf("server init error: %v", err)
	}
	defer serv.Close()

	if err := serv.RunServers(ctx); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
