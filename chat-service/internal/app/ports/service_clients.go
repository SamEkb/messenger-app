package ports

import (
	"context"

	friends "github.com/SamEkb/messenger-app/pkg/api/friends_service/v1"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
)

type UserServiceClient interface {
	GetUserProfile(ctx context.Context, userID string) (*UserProfile, error)
	GetProfiles(ctx context.Context, request *users.GetProfilesRequest) (*GetProfilesResponse, error)
	Close() error
}

type GetProfilesResponse struct {
	Profiles    map[string]*UserProfile
	NotFoundIds []string
}

type FriendServiceClient interface {
	CheckFriendsStatus(ctx context.Context, userID1, userID2 string) (friends.FriendshipStatus, error)
	CheckFriendshipsStatus(ctx context.Context, userIDs *friends.CheckFriendshipsStatusRequest) (*CheckFriendshipsStatusResponse, error)
	Close() error
}

type UserProfile struct {
	UserID      string
	Nickname    string
	Email       string
	Description string
	AvatarURL   string
}

type UserPair struct {
	UserID1 string
	UserID2 string
}

type CheckFriendshipsStatusResponse struct {
	NonFriendPairs []UserPair
	AllAreFriends  bool
}
