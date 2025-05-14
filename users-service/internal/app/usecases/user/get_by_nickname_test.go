package user

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/users-service/internal/app/usecases/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_GetByNickname(t *testing.T) {
	ctx := context.Background()

	testUUID := uuid.New()
	testUserID := models.UserID(testUUID)
	nickname := "testuser"

	type args struct {
		ctx      context.Context
		nickname string
	}

	tests := map[string]struct {
		args    args
		want    *ports.UserDto
		wantErr bool
		deps    func(t *testing.T) UseCase
	}{
		"get success": {
			args: args{
				ctx:      ctx,
				nickname: nickname,
			},
			want: &ports.UserDto{
				ID:          testUserID.String(),
				Email:       "test@test.ru",
				Nickname:    nickname,
				Description: "",
				AvatarURL:   "",
			},
			wantErr: false,
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))

				expectedUser, err := models.NewUser(testUserID, "test@test.ru", nickname, "", "")
				assert.NoError(t, err)

				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					GetByNickname(ctx, nickname).
					Return(expectedUser, nil).
					Once()

				return UseCase{
					userRepository: mockUserRepository,
					logger:         logger,
				}
			},
		},
		"failed to get user": {
			args: args{
				ctx:      ctx,
				nickname: nickname,
			},
			want:    nil,
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))

				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					GetByNickname(ctx, nickname).
					Return(nil, assert.AnError).
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
			dto, err := useCase.GetByNickname(tc.args.ctx, tc.args.nickname)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, dto)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, dto)
			}
		})
	}
}
