package ports

import (
	"context"

	auth "github.com/SamEkb/messenger-app/pkg/api/auth_service/v1"
	"github.com/SamEkb/messenger-app/pkg/api/events/v1"
)

type UserEventsKafkaProducer interface {
	ProduceUserRegisteredEvent(ctx context.Context, event *events.UserRegisteredEvent) error
}

type UserGrpcServer interface {
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
}
