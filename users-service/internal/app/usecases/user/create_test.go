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
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Create(t *testing.T) {
	ctx := context.Background()

	testUUID := uuid.New()
	testUserID := models.UserID(testUUID)

	type args struct {
		ctx context.Context
		dto *ports.UserDto
	}

	tests := map[string]struct {
		args    args
		want    string
		wantErr bool
		deps    func(t *testing.T) UseCase
	}{
		"create success": {
			args: args{
				ctx: ctx,
				dto: &ports.UserDto{
					ID:          testUUID.String(),
					Email:       "test@test.com",
					Nickname:    "testuser",
					Description: "Test user description",
				},
			},
			want:    testUUID.String(),
			wantErr: false,
			deps: func(t *testing.T) UseCase {
				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					Create(ctx, mock.AnythingOfType("*models.User")).
					Return(testUserID, nil)

				return UseCase{
					userRepository: mockUserRepository,
					logger:         logger.NewMockLogger(),
				}
			},
		},
		"failed to create": {
			args: args{
				ctx: ctx,
				dto: &ports.UserDto{
					ID:          testUUID.String(),
					Email:       "test@test.com",
					Nickname:    "testuser",
					Description: "Test user description",
				},
			},
			want:    "",
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				mockUserRepository := mocks.NewUserRepository(t)
				mockUserRepository.EXPECT().
					Create(ctx, mock.AnythingOfType("*models.User")).
					Return(models.UserID{}, assert.AnError)

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
			userID, err := useCase.Create(tc.args.ctx, tc.args.dto)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.want, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, userID)
			}
		})
	}
}
