package grpc

import (
	"context"
	"github.com/SamEkb/messenger-app/pkg/platform/middleware/resilience"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
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
	clientInterceptor := resilience.NewClientInterceptor(
		f.logger,
		f.config.RateLimit.DefaultLimit,
		f.config.RateLimit.DefaultBurst,
	)
	cbInterceptor := resilience.NewCircuitBreakerInterceptor(
		f.logger,
		resilience.WithFailureRatio(f.config.CircuitBreaker.FailureRatio),
		resilience.WithInterval(f.config.CircuitBreaker.Interval),
		resilience.WithTimeout(f.config.CircuitBreaker.Timeout),
		resilience.WithName(f.config.CircuitBreaker.Name),
		resilience.WithMaxRequests(f.config.CircuitBreaker.MaxRequests),
		resilience.WithMinRequests(f.config.CircuitBreaker.MinRequests),
		resilience.WithServerErrorCodes(f.config.CircuitBreaker.ServerErrorCodes),
	)
	retryInterceptor := resilience.RetryUnaryClientInterceptor(
		f.config.RetryConfig.MaxRetries,
		f.config.RetryConfig.RetryDelay,
		f.logger,
	)

	conn, err := grpc.DialContext(
		ctx,
		f.config.Users.Addr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			clientInterceptor,
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
