package clients

import (
	"context"
	"fmt"
	"time"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const friendsServiceAddr = "localhost:9003"

// FriendsClient provides a wrapper around the Friends Service gRPC client
type FriendsClient struct {
	client friends.FriendsServiceClient
	conn   *grpc.ClientConn
}

// NewFriendsClient creates a new Friends Service client
func NewFriendsClient() (*FriendsClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		friendsServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Friends Service: %w", err)
	}

	client := friends.NewFriendsServiceClient(conn)
	return &FriendsClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *FriendsClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CheckFriendshipStatus checks the friendship status between two users
func (c *FriendsClient) CheckFriendshipStatus(ctx context.Context, userID, friendID string) (*friends.CheckFriendshipStatusResponse, error) {
	req := &friends.CheckFriendshipStatusRequest{
		UserId:   userID,
		FriendId: friendID,
	}

	resp, err := c.client.CheckFriendshipStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check friendship status: %w", err)
	}

	return resp, nil
}
