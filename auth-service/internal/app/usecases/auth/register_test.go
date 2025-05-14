package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth/mocks"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Register(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		dto *ports.RegisterDto
	}

	tests := map[string]struct {
		args    args
		wantErr bool
		deps    func(t *testing.T) UseCase
	}{
		"register successful": {
			args: args{
				ctx: ctx,
				dto: &ports.RegisterDto{
					Username: "testuser",
					Email:    "test@test.ru",
					Password: "strongAndLongPassword",
				},
			},
			wantErr: false,
			deps: func(t *testing.T) UseCase {
				mockAuthRepo := mocks.NewAuthRepository(t)
				mockAuthRepo.EXPECT().
					Create(ctx, mock.AnythingOfType("*models.User")).
					Return(nil).
					Once()

				mockKafkaProducer := mocks.NewUserEventsKafkaProducer(t)
				mockKafkaProducer.EXPECT().
					ProduceUserRegisteredEvent(ctx, mock.AnythingOfType("*events.UserRegisteredEvent")).
					Return(nil).
					Once()

				return UseCase{
					authRepo:           mockAuthRepo,
					userEventPublisher: mockKafkaProducer,
					logger:             logger.NewLogger("local", "test"),
				}
			},
		},
		"failed to save user": {
			args: args{
				ctx: ctx,
				dto: &ports.RegisterDto{
					Username: "testuser",
					Email:    "test@test.ru",
					Password: "strongAndLongPassword",
				},
			},
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				mockAuthRepo := mocks.NewAuthRepository(t)
				mockAuthRepo.EXPECT().
					Create(ctx, mock.AnythingOfType("*models.User")).
					Return(errors.New("database error")).
					Once()

				return UseCase{
					authRepo: mockAuthRepo,
					logger:   logger.NewLogger("local", "test"),
				}
			},
		},
		"failed to publish event": {
			args: args{
				ctx: ctx,
				dto: &ports.RegisterDto{
					Username: "testuser",
					Email:    "test@test.ru",
					Password: "strongAndLongPassword",
				},
			},
			wantErr: true,
			deps: func(t *testing.T) UseCase {
				mockAuthRepo := mocks.NewAuthRepository(t)
				mockAuthRepo.EXPECT().
					Create(ctx, mock.AnythingOfType("*models.User")).
					Return(nil).
					Once()

				mockKafkaProducer := mocks.NewUserEventsKafkaProducer(t)
				mockKafkaProducer.EXPECT().
					ProduceUserRegisteredEvent(ctx, mock.AnythingOfType("*events.UserRegisteredEvent")).
					Return(errors.New("kafka error")).
					Once()

				return UseCase{
					authRepo:           mockAuthRepo,
					userEventPublisher: mockKafkaProducer,
					logger:             logger.NewLogger("local", "test"),
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			useCase := tc.deps(t)
			id, err := useCase.Register(tc.args.ctx, tc.args.dto)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, id)
			}
		})
	}
}
