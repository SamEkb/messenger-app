package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamEkb/messenger-app/chat-service/internal/app/ports"
	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"google.golang.org/grpc"
)

type UsersServiceClientAdapter struct {
	client users.UsersServiceClient
	conn   *grpc.ClientConn
}

func (c *UsersServiceClientAdapter) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *UsersServiceClientAdapter) GetUserProfile(userID string) (*ports.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &users.GetUserProfileRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUserProfile(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &ports.UserProfile{
		UserID:      userID,
		Nickname:    resp.Nickname,
		Email:       resp.Email,
		Description: resp.Description,
		AvatarURL:   resp.AvatarUrl,
	}, nil
}

func (c *UsersServiceClientAdapter) GetProfiles(ctx context.Context, request *users.GetProfilesRequest) (*ports.GetProfilesResponse, error) {
	resp, err := c.client.GetProfiles(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profiles: %w", err)
	}

	profiles := make(map[string]*ports.UserProfile)
	for id, profile := range resp.Profiles {
		profiles[id] = &ports.UserProfile{
			UserID:      profile.UserId,
			Nickname:    profile.Nickname,
			Email:       profile.Email,
			Description: profile.Description,
			AvatarURL:   profile.AvatarUrl,
		}
	}

	return &ports.GetProfilesResponse{
		Profiles:    profiles,
		NotFoundIds: resp.NotFoundIds,
	}, nil
}
