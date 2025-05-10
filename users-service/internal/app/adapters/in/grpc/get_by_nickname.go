package grpc

import (
	"context"

	users "github.com/SamEkb/messenger-app/pkg/api/users_service/v1"
)

func (s *UsersServiceServer) GetUserProfileByNickname(ctx context.Context, req *users.GetUserProfileByNicknameRequest) (*users.GetUserProfileByNicknameResponse, error) {
	s.logger.Info("Getting user profile by nickname")

	user, err := s.userUseCase.GetByNickname(ctx, req.GetNickname())
	if err != nil {
		s.logger.Error("Failed to get user profile by nickname", "error", err)
		return nil, err
	}

	s.logger.Info("User profile successfully retrieved")

	return &users.GetUserProfileByNicknameResponse{
		Nickname:    user.Nickname,
		Email:       user.Email,
		Description: user.Description,
		AvatarUrl:   user.AvatarURL,
	}, nil
}
