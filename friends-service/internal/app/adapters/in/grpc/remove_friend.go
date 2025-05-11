package grpc

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
)

func (s *FriendshipServiceServer) RemoveFriend(ctx context.Context, req *friends.RemoveFriendRequest) (*friends.RemoveFriendResponse, error) {
	s.logger.Info("removing friend")

	if err := s.friendshipUseCase.DeleteFriend(ctx, req.UserId, req.GetFriendId()); err != nil {
		s.logger.Error("failed to remove friend", "error", err)
		return nil, err
	}

	s.logger.Info("friend removed")

	return &friends.RemoveFriendResponse{
		Success: true,
		Message: "Friend removed",
	}, nil
}
