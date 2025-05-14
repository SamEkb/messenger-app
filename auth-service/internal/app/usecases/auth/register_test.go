package auth

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/usecases/auth/mocks"
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
		err     error
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
			err:     nil,
			deps: func(t *testing.T) UseCase {
				var buf bytes.Buffer
				logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))

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
					logger:             logger,
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
				assert.IsType(t, tc.err, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, models.UserID{}, id)
				assert.NotEmpty(t, id.String())
			}
		})
	}
}
