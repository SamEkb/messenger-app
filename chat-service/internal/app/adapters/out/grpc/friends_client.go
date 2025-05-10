package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"google.golang.org/grpc"
)

type FriendsServiceClientAdapter struct {
	client friends.FriendsServiceClient
	conn   *grpc.ClientConn
}

func (c *FriendsServiceClientAdapter) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *FriendsServiceClientAdapter) CheckFriendsStatus(userID1, userID2 string) (friends.FriendshipStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &friends.CheckFriendshipStatusRequest{
		UserId:   userID1,
		FriendId: userID2,
	}

	resp, err := c.client.CheckFriendshipStatus(ctx, req)
	if err != nil {
		return friends.FriendshipStatus_FRIENDSHIP_STATUS_UNSPECIFIED,
			fmt.Errorf("failed to check friendship status: %w", err)
	}

	return resp.Status, nil
}

func (c *FriendsServiceClientAdapter) CheckFriendshipsStatus(ctx context.Context, request *friends.CheckFriendshipsStatusRequest) (*ports.CheckFriendshipsStatusResponse, error) {
	resp, err := c.client.CheckFriendshipsStatus(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to check friendships status: %w", err)
	}

	userPairs := make([]ports.UserPair, 0, len(resp.NonFriendPairs))
	for _, pair := range resp.NonFriendPairs {
		userPairs = append(userPairs, ports.UserPair{
			UserID1: pair.UserId1,
			UserID2: pair.UserId2,
		})
	}

	return &ports.CheckFriendshipsStatusResponse{
		NonFriendPairs: userPairs,
		AllAreFriends:  resp.AllAreFriends,
	}, nil
}
