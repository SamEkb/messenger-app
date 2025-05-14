package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth/mocks"
	customerrors "github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_Login(t *testing.T) {
	ttlDuration := time.Hour
	ctx := context.Background()

	userID := models.UserID(uuid.New())
	email := "test@test.ru"
	password := "somestrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	type args struct {
		ctx context.Context
		dto *ports.LoginDto
	}

	tests := map[string]struct {
		args        args
		want        models.Token
		wantErr     bool
		expectedErr error
		deps        func(t *testing.T) UseCase
	}{
		"login successful": {
			args: args{
				ctx: ctx,
				dto: &ports.LoginDto{
					Email:    email,
					Password: password,
				},
			},
			want:    models.Token("very-strong-token"),
			wantErr: false,
			deps: func(t *testing.T) UseCase {
				logger := logger.NewLogger("local", "test")

				mockUser, err := models.NewUser(userID, "testuser", email, hashedPassword)
				assert.NoError(t, err)

				mockAuthRepo := mocks.NewAuthRepository(t)
				mockAuthRepo.EXPECT().
					FindUserByEmail(ctx, email).
					Return(mockUser, nil).
					Once()

				expectedToken := models.Token("very-strong-token")
				mockAuthToken, err := models.NewAuthToken(
					expectedToken,
					userID,
					time.Now().Add(ttlDuration),
				)
				assert.NoError(t, err)

				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().
					Create(ctx, mock.AnythingOfType("*models.AuthToken")).
					Return(mockAuthToken, nil).
					Once()

				return UseCase{
					authRepo:  mockAuthRepo,
					tokenRepo: mockTokenRepo,
					logger:    logger,
					tokenTTL:  ttlDuration,
				}
			},
		},
		"user not found": {
			args: args{
				ctx: ctx,
				dto: &ports.LoginDto{
					Email:    "test@test.ru",
					Password: password,
				},
			},
			want:    "",
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				logger := logger.NewLogger("local", "test")

				mockAuthRepo := mocks.NewAuthRepository(t)
				mockAuthRepo.EXPECT().
					FindUserByEmail(ctx, "test@test.ru").
					Return(nil, errors.New("user not found")).
					Once()

				mockTokenRepo := mocks.NewTokenRepository(t)

				return UseCase{
					authRepo:  mockAuthRepo,
					tokenRepo: mockTokenRepo,
					logger:    logger,
					tokenTTL:  ttlDuration,
				}
			},
		},
		"invalid credentials": {
			args: args{
				ctx: ctx,
				dto: &ports.LoginDto{
					Email:    email,
					Password: "wrongPassword",
				},
			},
			want:        "",
			wantErr:     true,
			expectedErr: customerrors.NewUnauthorizedError("invalid credentials"),
			deps: func(t *testing.T) UseCase {
				logger := logger.NewLogger("local", "test")

				mockUser, err := models.NewUser(userID, "testuser", email, hashedPassword)
				assert.NoError(t, err)

				mockAuthRepo := mocks.NewAuthRepository(t)
				mockAuthRepo.EXPECT().
					FindUserByEmail(ctx, email).
					Return(mockUser, nil).
					Once()
				mockTokenRepo := mocks.NewTokenRepository(t)

				return UseCase{
					authRepo:  mockAuthRepo,
					tokenRepo: mockTokenRepo,
					logger:    logger,
					tokenTTL:  ttlDuration,
				}
			},
		},
		"failed to save token": {
			args: args{
				ctx: ctx,
				dto: &ports.LoginDto{
					Email:    email,
					Password: password,
				},
			},
			want:        "",
			wantErr:     true,
			expectedErr: customerrors.NewInternalError(nil, "failed to save token"),
			deps: func(t *testing.T) UseCase {
				logger := logger.NewLogger("local", "test")

				mockUser, err := models.NewUser(userID, "testuser", email, hashedPassword)
				assert.NoError(t, err)

				mockAuthRepo := mocks.NewAuthRepository(t)
				mockAuthRepo.EXPECT().
					FindUserByEmail(ctx, email).
					Return(mockUser, nil).
					Once()

				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().
					Create(ctx, mock.AnythingOfType("*models.AuthToken")).
					Return(nil, customerrors.NewInternalError(nil, "database error")).
					Once()

				return UseCase{
					authRepo:  mockAuthRepo,
					tokenRepo: mockTokenRepo,
					logger:    logger,
					tokenTTL:  ttlDuration,
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			useCase := tc.deps(t)
			token, err := useCase.Login(tc.args.ctx, tc.args.dto)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.IsType(t, tc.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, token)
			}
		})
	}
}
