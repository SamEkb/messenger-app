package grpc

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
)

func (s *FriendshipServiceServer) SendFriendRequest(ctx context.Context, req *friends.SendFriendRequestRequest) (*friends.SendFriendRequestResponse, error) {
	s.logger.Info("sending friend request")

	if err := s.friendshipUseCase.SendFriendRequest(ctx, req.GetUserId(), req.GetFriendId()); err != nil {
		s.logger.Error("failed to send friend request", "error", err)
		return nil, err
	}

	return &friends.SendFriendRequestResponse{
		Success: true,
		Message: "Friend request sent",
	}, nil
}
