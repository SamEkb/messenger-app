package user

import (
	"context"
	"testing"

	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/users-service/internal/app/usecases/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_Get(t *testing.T) {
	ctx := context.Background()

	testUUID := uuid.New()
	testUserID := models.UserID(testUUID)

	type args struct {
		ctx context.Context
		id  string
	}

	tests := map[string]struct {
		args    args
		want    *ports.UserDto
		wantErr bool
		deps    func(t *testing.T) UseCase
	}{
		"get success": {
			args: args{
				ctx: ctx,
				id:  testUUID.String(),
			},
			want: &ports.UserDto{
				ID:          testUserID.String(),
				Email:       "test@test.ru",
				Nickname:    "testuser",
				Description: "",
				AvatarURL:   "",
			},
			wantErr: false,
			deps: func(t *testing.T) UseCase {
				expectedUser, err := models.NewUser(testUserID, "test@test.ru", "testuser", "", "")
				assert.NoError(t, err)

				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					Get(ctx, testUserID).
					Return(expectedUser, nil).
					Once()

				return UseCase{
					userRepository: mockUserRepository,
					logger:         logger.NewMockLogger(),
				}
			},
		},
		"failed to get user": {
			args: args{
				ctx: ctx,
				id:  testUUID.String(),
			},
			want:    nil,
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					Get(ctx, testUserID).
					Return(nil, assert.AnError).
					Once()

				return UseCase{
					userRepository: mockUserRepository,
					logger:         logger.NewMockLogger(),
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			useCase := tc.deps(t)
			dto, err := useCase.Get(tc.args.ctx, tc.args.id)

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
