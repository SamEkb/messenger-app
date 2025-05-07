package auth

import (
	"context"
	"errors"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a *UseCase) Register(ctx context.Context, dto *ports.RegisterDto) (*models.User, error) {
	user, err := a.authRepo.FindUserByEmail(ctx, dto.Email)
	if err == nil {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New()
	newUser, err := models.NewUser(models.UserID(id), dto.Username, dto.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	if err = a.authRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
