package main

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/repositories/auth/in_memory"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repository := in_memory.NewRepository()

}
