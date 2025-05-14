package user

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/users-service/internal/app/usecases/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Update(t *testing.T) {
	ctx := context.Background()

	testUUID := uuid.New()

	type args struct {
		ctx context.Context
		dto *ports.UserDto
	}

	tests := map[string]struct {
		args    args
		wantErr bool
		deps    func(t *testing.T) UseCase
	}{
		"update success": {
			args: args{
				ctx: ctx,
				dto: &ports.UserDto{
					ID:          testUUID.String(),
					Email:       "updated@test.com",
					Nickname:    "updateduser",
					Description: "Test user description",
				},
			},
			wantErr: false,
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))

				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					Update(ctx, mock.AnythingOfType("*models.User")).
					Return(nil).
					Once()

				return UseCase{
					userRepository: mockUserRepository,
					logger:         logger,
				}
			},
		},
		"failed to update": {
			args: args{
				ctx: ctx,
				dto: &ports.UserDto{
					ID:          testUUID.String(),
					Email:       "updated@test.com",
					Nickname:    "updateduser",
					Description: "Test user description",
				},
			},
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))

				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					Update(ctx, mock.AnythingOfType("*models.User")).
					Return(assert.AnError).
					Once()

				return UseCase{
					userRepository: mockUserRepository,
					logger:         logger,
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			useCase := tc.deps(t)
			err := useCase.Update(tc.args.ctx, tc.args.dto)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
