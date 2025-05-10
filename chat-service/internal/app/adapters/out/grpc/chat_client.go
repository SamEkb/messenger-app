package grpc

import (
	"context"
	"log/slog"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/chat-service/pkg/errors"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
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

func (f *Client) NewFriendsServiceClient(ctx context.Context) (ports.FriendServiceClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		f.config.Friends.Addr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.NewServiceError(err, "failed to connect to Friends Service")
	}

	client := friends.NewFriendsServiceClient(conn)
	return &FriendsServiceClientAdapter{
		client: client,
		conn:   conn,
	}, nil
}
