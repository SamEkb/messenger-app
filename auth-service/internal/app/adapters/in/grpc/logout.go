package grpc

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/pkg/errors"
	auth "github.com/SamEkb/messenger-app/pkg/api/auth_service/v1"
)

func (s *Server) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		s.logger.Error("validation error", "error", err)
		return nil, errors.NewValidationError("invalid request: %v", err).
			WithDetails("token", req.GetToken())
	}

	token := models.Token(req.GetToken())

	s.logger.Info("logout request received", "token", token)

	if err := s.authUseCase.Logout(ctx, token); err != nil {
		s.logger.Error("logout failed", "error", err)
		return nil, err
	}

	s.logger.Info("logout successful")

	return &auth.LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}
