package grpc

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *FriendshipServiceServer) CheckFriendshipStatus(ctx context.Context, req *friends.CheckFriendshipStatusRequest) (*friends.CheckFriendshipStatusResponse, error) {
	s.logger.Info("checking friendship status", "user_id", req.GetUserId(), "friend_id", req.GetFriendId())

	friendsList, err := s.friendshipUseCase.GetFriends(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("failed to get friends list", "error", err)
		return nil, err
	}

	var status friends.FriendshipStatus
	var createdAt, updatedAt *timestamppb.Timestamp

	status = friends.FriendshipStatus_FRIENDSHIP_STATUS_UNSPECIFIED

	for _, friendship := range friendsList {
		if friendship.RecipientID() == req.GetFriendId() || friendship.RequestorID() == req.GetFriendId() {
			status = mapStatusToProto(friendship.Status())
			createdAt = timestamppb.New(friendship.CreatedAt())
			updatedAt = timestamppb.New(friendship.UpdatedAt())
			break
		}
	}

	return &friends.CheckFriendshipStatusResponse{
		Status:    status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
