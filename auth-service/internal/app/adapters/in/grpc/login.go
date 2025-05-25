package grpc

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	auth "github.com/SamEkb/messenger-app/pkg/api/auth_service/v1"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
)

func (s *Server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error("validation error", "error", err)
		return nil, errors.NewValidationError("invalid request: %v", err).
			WithDetails("email", req.GetEmail())
	}

	s.logger.Info("login request received", "email", req.GetEmail())

	loginDTO := &ports.LoginDto{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	token, err := s.authUseCase.Login(ctx, loginDTO)
	if err != nil {
		s.logger.Error("login failed", "error", err)
		return nil, err
	}

	s.logger.Info("login successful", "token", token)

	return &auth.LoginResponse{
		Token:     string(token),
		UserId:    "",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Success:   true,
		Message:   "Login successful",
	}, nil
}
