package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth/mocks"
	customerrors "github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
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
		args        args
		wantErr     bool
		expectedErr error
		deps        func(t *testing.T) UseCase
	}{
		"logout successful": {
			args: args{
				ctx:   ctx,
				token: models.Token("very-strong-token"),
			},
			wantErr: false,
			deps: func(t *testing.T) UseCase {
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
					logger:    logger.NewMockLogger(),
				}
			},
		},
		"failed to validate token": {
			args: args{
				ctx:   ctx,
				token: models.Token("very-strong-token"),
			},
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().ValidateToken(ctx, "very-strong-token").
					Return(false, errors.New("database error")).
					Once()

				return UseCase{
					tokenRepo: mockTokenRepo,
					logger:    logger.NewMockLogger(),
				}
			},
		},
		"invalid token": {
			args: args{
				ctx:   ctx,
				token: models.Token("invalid-token"),
			},
			wantErr:     true,
			expectedErr: customerrors.NewTokenError(nil, "invalid token"),
			deps: func(t *testing.T) UseCase {
				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().ValidateToken(ctx, "invalid-token").
					Return(false, nil).
					Once()

				return UseCase{
					tokenRepo: mockTokenRepo,
					logger:    logger.NewMockLogger(),
				}
			},
		},
		"failed to delete token": {
			args: args{
				ctx:   ctx,
				token: models.Token("very-strong-token"),
			},
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				mockTokenRepo := mocks.NewTokenRepository(t)
				mockTokenRepo.EXPECT().ValidateToken(ctx, "very-strong-token").
					Return(true, nil).
					Once()
				mockTokenRepo.EXPECT().
					DeleteToken(ctx, mock.AnythingOfType("models.Token")).
					Return(errors.New("database error")).
					Once()

				return UseCase{
					tokenRepo: mockTokenRepo,
					logger:    logger.NewMockLogger(),
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
				if tc.expectedErr != nil {
					assert.IsType(t, tc.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
