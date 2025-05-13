package grpc

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
)

func (s *FriendshipServiceServer) AcceptFriendRequest(ctx context.Context, req *friends.AcceptFriendRequestRequest) (*friends.AcceptFriendRequestResponse, error) {
	s.logger.Info("accepting friend request")

	if err := s.friendshipUseCase.AcceptFriendRequest(ctx, req.GetUserId(), req.GetFriendId()); err != nil {
		s.logger.Error("failed to accept friend request", "error", err)
		return nil, err
	}

	s.logger.Info("friend request accepted")

	return &friends.AcceptFriendRequestResponse{
		Success: true,
		Message: "Friend request accepted",
	}, nil
}
