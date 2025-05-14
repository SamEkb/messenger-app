package auth

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth/mocks"
	"github.com/SamEkb/messenger-app/auth-service/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Logout(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		token models.Token
	}
	tests := map[string]struct {
		args    args
		wantErr bool
		err     error
		deps    func(t *testing.T) UseCase
	}{
		"logout successful": {
			args: args{
				ctx:   ctx,
				token: models.Token("very-strong-token"),
			},
			wantErr: false,
			err:     nil,
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))

				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().ValidateToken(ctx, "very-strong-token").
					Return(true, nil).
					Once()

				mockTokenRepo.EXPECT().
					DeleteToken(ctx, mock.AnythingOfType("models.Token")).
					Return(nil).
					Once()

				return UseCase{
					tokenRepo: mockTokenRepo,
					logger:    logger,
				}
			},
		},
		"failed to validate token": {
			args: args{
				ctx:   ctx,
				token: models.Token("very-strong-token"),
			},
			wantErr: true,
			err:     errors.NewInternalError(nil, "failed to validate token"),
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().ValidateToken(ctx, "very-strong-token").
					Return(false, errors.NewInternalError(nil, "database error")).
					Once()

				return UseCase{
					tokenRepo: mockTokenRepo,
					logger:    logger,
				}
			},
		},
		"invalid token": {
			args: args{
				ctx:   ctx,
				token: models.Token("invalid-token"),
			},
			wantErr: true,
			err:     errors.NewTokenError(errors.ErrInvalidToken, "token is invalid"),
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().ValidateToken(ctx, "invalid-token").
					Return(false, nil).
					Once()

				return UseCase{
					tokenRepo: mockTokenRepo,
					logger:    logger,
				}
			},
		},
		"failed to delete token": {
			args: args{
				ctx:   ctx,
				token: models.Token("very-strong-token"),
			},
			wantErr: true,
			err:     errors.NewInternalError(nil, "failed to delete token"),
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().ValidateToken(ctx, "very-strong-token").
					Return(true, nil).
					Once()
				mockTokenRepo.EXPECT().
					DeleteToken(ctx, mock.AnythingOfType("models.Token")).
					Return(errors.NewInternalError(nil, "database error")).
					Once()

				return UseCase{
					tokenRepo: mockTokenRepo,
					logger:    logger,
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			useCase := tc.deps(t)
			err := useCase.Logout(tc.args.ctx, tc.args.token)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
