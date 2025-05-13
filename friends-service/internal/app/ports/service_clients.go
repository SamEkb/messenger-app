package ports

import (
	"context"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
)

type UserServiceClient interface {
	GetUserProfile(userID string) (*UserProfile, error)
	GetProfiles(ctx context.Context, request *users.GetProfilesRequest) (*GetProfilesResponse, error)
}

type GetProfilesResponse struct {
	Profiles    map[string]*UserProfile
	NotFoundIds []string
}

type UserProfile struct {
	UserID      string
	Nickname    string
	Email       string
	Description string
	AvatarURL   string
}
