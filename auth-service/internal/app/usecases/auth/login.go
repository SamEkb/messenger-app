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
	a.logger.Debug("login attempt", "email", dto.Email)

	user, err := a.authRepo.FindUserByEmail(ctx, dto.Email)
	if err != nil {
		a.logger.Warn("user not found during login", "email", dto.Email, "error", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(dto.Password)); err != nil {
		a.logger.Warn("invalid credentials", "email", dto.Email)
		return "", errors.New("invalid credentials")
	}

	token, err := models.NewAuthToken(
		models.Token(uuid.New().String()),
		user.ID(),
		time.Now().Add(a.tokenTTL),
	)
	if err != nil {
		a.logger.Error("failed to create auth token", "error", err)
		return "", err
	}

	token, err = a.tokenRepo.Create(ctx, token)
	if err != nil {
		a.logger.Error("failed to save token", "error", err)
		return "", err
	}

	a.logger.Info("login successful", "user_id", user.ID(), "token", token.Token())
	return token.Token(), nil
}
