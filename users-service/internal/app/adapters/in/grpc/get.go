package grpc

import (
	"context"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
)

func (s *UsersServiceServer) GetUserProfile(ctx context.Context, req *users.GetUserProfileRequest) (*users.GetUserProfileResponse, error) {
	s.logger.Info("Getting user profile")

	user, err := s.userUseCase.Get(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("Failed to get user profile", "error", err)
		return nil, err
	}

	s.logger.Info("User profile successfully retrieved")

	return &users.GetUserProfileResponse{
		Nickname:    user.Nickname,
		Email:       user.Email,
		Description: user.Description,
		AvatarUrl:   user.AvatarURL,
	}, nil
}
