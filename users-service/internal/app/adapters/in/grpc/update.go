package grpc

import (
	"context"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

func (s *UsersServiceServer) UpdateUserProfile(ctx context.Context, req *users.UpdateUserProfileRequest) (*users.UpdateUserProfileResponse, error) {
	s.logger.Info("Updating user profile")

	dto := &ports.UserDto{
		ID:          req.UserId,
		Email:       req.Email,
		Nickname:    req.Nickname,
		Description: req.Description,
		AvatarURL:   req.AvatarUrl,
	}
	if err := s.userUseCase.Update(ctx, dto); err != nil {
		s.logger.Error("Failed to update user profile", "error", err)
		return nil, err
	}

	return &users.UpdateUserProfileResponse{
		Success: true,
		Message: "User profile updated successfully",
	}, nil
}
