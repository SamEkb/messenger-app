package grpc

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
)

func (s *FriendshipServiceServer) RejectFriendRequest(ctx context.Context, req *friends.RejectFriendRequestRequest) (*friends.RejectFriendRequestResponse, error) {
	s.logger.Info("rejecting friend request")

	if err := s.friendshipUseCase.RejectFriendRequest(ctx, req.GetUserId(), req.GetFriendId()); err != nil {
		s.logger.Error("failed to reject friend request", "error", err)
		return nil, err
	}

	s.logger.Info("friend request rejected")

	return &friends.RejectFriendRequestResponse{
		Success: true,
		Message: "Friend request rejected",
	}, nil
}
