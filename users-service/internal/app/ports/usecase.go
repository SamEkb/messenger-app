package ports

import (
	"context"
)

type UserUseCase interface {
	Create(ctx context.Context, dto *UserDto) (string, error)
	Get(ctx context.Context, id string) (*UserDto, error)
	GetByNickname(ctx context.Context, nickname string) (*UserDto, error)
	Update(ctx context.Context, dto *UserDto) error
}

type UserDto struct {
	ID          string
	Email       string
	Nickname    string
	Description string
	AvatarURL   string
}
