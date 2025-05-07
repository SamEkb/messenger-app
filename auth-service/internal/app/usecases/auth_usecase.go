package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo  ports.AuthRepository
	tokenRepo ports.TokenRepository
	tokenTTL  time.Duration
}

func NewAuthUsecase(authRepo ports.AuthRepository, tokenRepo ports.TokenRepository, tokenTTL time.Duration) *AuthUsecase {
	return &AuthUsecase{
		authRepo:  authRepo,
		tokenRepo: tokenRepo,
		tokenTTL:  tokenTTL,
	}
}

func (a *AuthUsecase) Login(ctx context.Context, dto *ports.LoginDto) (*models.Token, error) {
	user, err := a.authRepo.FindUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(dto.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := models.NewAuthToken(
		models.Token(uuid.New().String()),
		user.ID(),
		time.Now().Add(a.tokenTTL),
	)
	if err != nil {
		return nil, err
	}

	token, err = a.tokenRepo.Create(ctx, token)
	if err != nil {
		return nil, err
	}

	return &token.Token(), nil
}

func (a *AuthUsecase) Register(ctx context.Context, dto *ports.RegisterDto) (*models.User, error) {
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

func (a *AuthUsecase) Logout(ctx context.Context, token models.Token) error {
	a.tokenRepo.DeleteToken(ctx, token)
	return nil
}
