package auth

import (
	"context"
	"errors"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a *UseCase) Login(ctx context.Context, dto *ports.LoginDto) (models.Token, error) {
	user, err := a.authRepo.FindUserByEmail(ctx, dto.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(dto.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := models.NewAuthToken(
		models.Token(uuid.New().String()),
		user.ID(),
		time.Now().Add(a.tokenTTL),
	)
	if err != nil {
		return "", err
	}

	token, err = a.tokenRepo.Create(ctx, token)
	if err != nil {
		return "", err
	}

	return token.Token(), nil
}
