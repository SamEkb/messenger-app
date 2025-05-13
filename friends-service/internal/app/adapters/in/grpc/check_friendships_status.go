package grpc

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
)

func (s *FriendshipServiceServer) CheckFriendshipsStatus(ctx context.Context, req *friends.CheckFriendshipsStatusRequest) (*friends.CheckFriendshipsStatusResponse, error) {
	s.logger.Info("checking multiple friendships status", "user_count", len(req.GetUserIds()))

	nonFriendPairs, err := s.friendshipUseCase.CheckMultipleFriendships(ctx, req.GetUserIds())
	if err != nil {
		s.logger.Error("failed to check multiple friendships", "error", err)
		return nil, err
	}

	protoPairs := make([]*friends.CheckFriendshipsStatusResponse_UserPair, 0, len(nonFriendPairs))
	for _, pair := range nonFriendPairs {
		protoPairs = append(protoPairs, &friends.CheckFriendshipsStatusResponse_UserPair{
			UserId1: pair.UserID1,
			UserId2: pair.UserID2,
		})
	}

	allAreFriends := len(nonFriendPairs) == 0

	return &friends.CheckFriendshipsStatusResponse{
		NonFriendPairs: protoPairs,
		AllAreFriends:  allAreFriends,
	}, nil
}
