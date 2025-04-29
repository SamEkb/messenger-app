package clients

import (
	"context"
	"fmt"
	"time"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const usersServiceAddr = "localhost:9004"

// UsersClient provides a wrapper around the Users Service gRPC client
type UsersClient struct {
	client users.UsersServiceClient
	conn   *grpc.ClientConn
}

// NewUsersClient creates a new Users Service client
func NewUsersClient() (*UsersClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		usersServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Users Service: %w", err)
	}

	client := users.NewUsersServiceClient(conn)
	return &UsersClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *UsersClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetUserProfile retrieves a user profile by ID
func (c *UsersClient) GetUserProfile(ctx context.Context, userID string) (*users.GetUserProfileResponse, error) {
	req := &users.GetUserProfileRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUserProfile(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return resp, nil
}

// GetUserProfileByNickname retrieves a user profile by nickname
func (c *UsersClient) GetUserProfileByNickname(ctx context.Context, nickname string) (*users.GetUserProfileByNicknameResponse, error) {
	req := &users.GetUserProfileByNicknameRequest{
		Nickname: nickname,
	}

	resp, err := c.client.GetUserProfileByNickname(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile by nickname: %w", err)
	}

	return resp, nil
}
