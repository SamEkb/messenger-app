package grpc

import (
	"context"

	"github.com/SamEkb/messenger-app/chat-service/config/env"
	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	mw "github.com/SamEkb/messenger-app/pkg/platform/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config *env.ClientsConfig
	logger logger.Logger
}

func NewClient(config *env.ClientsConfig, logger logger.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

func (f *Client) NewUsersServiceClient(ctx context.Context) (ports.UserServiceClient, error) {
	cbInterceptor := mw.NewCircuitBreakerInterceptor(
		f.logger,
		mw.WithFailureRatio(f.config.CircuitBreaker.FailureRatio),
		mw.WithInterval(f.config.CircuitBreaker.Interval),
		mw.WithTimeout(f.config.CircuitBreaker.Timeout),
		mw.WithName(f.config.CircuitBreaker.Name),
		mw.WithMaxRequests(f.config.CircuitBreaker.MaxRequests),
		mw.WithMinRequests(f.config.CircuitBreaker.MinRequests),
		mw.WithServerErrorCodes(f.config.CircuitBreaker.ServerErrorCodes),
	)
	retryInterceptor := mw.RetryUnaryClientInterceptor(f.config.MaxRetries, f.config.RetryDelay, f.logger)

	conn, err := grpc.DialContext(
		ctx,
		f.config.Users.Addr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			cbInterceptor,
			retryInterceptor,
		),
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
	cbInterceptor := mw.NewCircuitBreakerInterceptor(
		f.logger,
		mw.WithFailureRatio(f.config.CircuitBreaker.FailureRatio),
		mw.WithInterval(f.config.CircuitBreaker.Interval),
		mw.WithTimeout(f.config.CircuitBreaker.Timeout),
		mw.WithName(f.config.CircuitBreaker.Name),
		mw.WithMaxRequests(f.config.CircuitBreaker.MaxRequests),
		mw.WithMinRequests(f.config.CircuitBreaker.MinRequests),
		mw.WithServerErrorCodes(f.config.CircuitBreaker.ServerErrorCodes),
	)
	retryInterceptor := mw.RetryUnaryClientInterceptor(f.config.MaxRetries, f.config.RetryDelay, f.logger)

	conn, err := grpc.DialContext(
		ctx,
		f.config.Friends.Addr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			cbInterceptor,
			retryInterceptor,
		),
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
