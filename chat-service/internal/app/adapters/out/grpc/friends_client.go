package grpc

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/chat-service/pkg/errors"
	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"google.golang.org/grpc"
	grpcStatus "google.golang.org/grpc/status"
)

type FriendsServiceClientAdapter struct {
	client friends.FriendsServiceClient
	conn   *grpc.ClientConn
}

func (c *FriendsServiceClientAdapter) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return errors.NewServiceError(err, "failed to close connection to Friends Service")
		}
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
		st, ok := grpcStatus.FromError(err)
		if ok {
			return resp.Status, errors.NewServiceError(err, "failed to check friendship status: %s", st.Message())
		}
		return resp.Status, errors.NewServiceError(err, "failed to check friendship status")
	}

	return resp.Status, nil
}

func (c *FriendsServiceClientAdapter) CheckFriendshipsStatus(ctx context.Context, request *friends.CheckFriendshipsStatusRequest) (*ports.CheckFriendshipsStatusResponse, error) {
	resp, err := c.client.CheckFriendshipsStatus(ctx, request)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			return nil, errors.NewServiceError(err, "failed to check friendships status: %s", st.Message())
		}
		return nil, errors.NewServiceError(err, "failed to check friendships status")
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
