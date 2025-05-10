package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientFactory struct {
	config *env.Config
	logger *slog.Logger
}

func NewClientFactory(config *env.Config, logger *slog.Logger) *ClientFactory {
	return &ClientFactory{
		config: config,
		logger: logger,
	}
}

func (f *ClientFactory) NewUsersServiceClient() (ports.UserServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.config.Clients.GRPC.ConnectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		f.config.Clients.Users.Addr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Users Service: %w", err)
	}

	client := users.NewUsersServiceClient(conn)
	return &UsersServiceClientAdapter{
		client: client,
		conn:   conn,
	}, nil
}

func (f *ClientFactory) NewFriendsServiceClient() (ports.FriendServiceClient, error) {
	var (
		conn   *grpc.ClientConn
		err    error
		retry  = 0
		logger = f.logger.With("service", "friends_client")
	)

	for retry < f.config.Clients.GRPC.RetryAttempts {
		ctx, cancel := context.WithTimeout(context.Background(), f.config.Clients.GRPC.ConnectionTimeout)
		conn, err = grpc.DialContext(
			ctx,
			f.config.Clients.Friends.Addr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		cancel()

		if err == nil {
			break
		}

		retry++
		logger.Warn("failed to connect to Friends Service, retrying...",
			"attempt", retry,
			"max_attempts", f.config.Clients.GRPC.RetryAttempts,
			"error", err)

		if retry < f.config.Clients.GRPC.RetryAttempts {
			time.Sleep(f.config.Clients.GRPC.RetryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Friends Service after %d attempts: %w",
			f.config.Clients.GRPC.RetryAttempts, err)
	}

	client := friends.NewFriendsServiceClient(conn)
	return &FriendsServiceClientAdapter{
		client: client,
		conn:   conn,
	}, nil
}
