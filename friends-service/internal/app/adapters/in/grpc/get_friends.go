package grpc

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *FriendshipServiceServer) GetFriendsList(ctx context.Context, req *friends.GetFriendsListRequest) (*friends.GetFriendsListResponse, error) {
	s.logger.Info("getting friends list")

	friendsList, err := s.friendshipUseCase.GetFriends(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("failed to get friends list", "error", err)
		return nil, err
	}

	protoFriends := make([]*friends.FriendInfo, 0, len(friendsList))
	for _, f := range friendsList {
		protoFriends = append(protoFriends, &friends.FriendInfo{
			UserId:    f.RecipientID(), // или f.requestorID, если нужен инициатор
			Nickname:  "",              // TODO: получить из users-service
			AvatarUrl: "",              // TODO: получить из users-service
			Status:    mapStatusToProto(f.status),
			CreatedAt: timestamppb.New(f.createdAt),
			UpdatedAt: timestamppb.New(f.updatedAt),
		})
	}

	return &friends.GetFriendsListResponse{
		Friends: protoFriends,
	}, nil
}

// Маппинг строки статуса в proto enum
func mapStatusToProto(status string) friends.FriendshipStatus {
	switch status {
	case "REQUESTED":
		return friends.FriendshipStatus_FRIENDSHIP_STATUS_REQUESTED
	case "ACCEPTED":
		return friends.FriendshipStatus_FRIENDSHIP_STATUS_ACCEPTED
	case "REJECTED":
		return friends.FriendshipStatus_FRIENDSHIP_STATUS_REJECTED
	default:
		return friends.FriendshipStatus_FRIENDSHIP_STATUS_UNSPECIFIED
	}
}
