package grpc

import (
	"context"
	"log/slog"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/friends-service/pkg/errors"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config *env.ClientsConfig
	logger *slog.Logger
}

func NewClient(config *env.ClientsConfig, logger *slog.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

func (f *Client) NewUsersServiceClient(ctx context.Context) (ports.UserServiceClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		f.config.Users.Addr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.NewServiceError(err, "failed to connect to Users Service")
	}

	client := users.NewUsersServiceClient(conn)
	return &UsersServiceClientAdapter{
		client: client,
		conn:   conn,
	}, nil
}
