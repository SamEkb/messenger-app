package grpc

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	auth "github.com/SamEkb/messenger-app/pkg/api/auth_service/v1"
)

func (s *Server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error("validation error", "error", err)
		return nil, err
	}

	s.logger.Info("register request received",
		"username", req.GetUsername(),
		"email", req.GetEmail())

	registerDTO := &ports.RegisterDto{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	userID, err := s.authUseCase.Register(ctx, registerDTO)
	if err != nil {
		s.logger.Error("failed to register user", "error", err)
		return nil, err
	}

	s.logger.Info("user registered successfully", "userID", userID)

	return &auth.RegisterResponse{
		UserId:  userID.String(),
		Message: "User registered successfully",
		Success: true,
	}, nil
}
