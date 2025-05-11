package in_memory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/users-service/pkg/errors"
)

var _ ports.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	mu    sync.RWMutex
	users map[models.UserID]*models.User

	logger *slog.Logger
}

func NewUserRepository(logger *slog.Logger) *UserRepository {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(nil, nil))
	}

	return &UserRepository{
		users:  make(map[models.UserID]*models.User),
		logger: logger.With("component", "user_repository"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (models.UserID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debug("attempting to create new user", "user_id", user.ID(), "email", user.Email())

	if _, ok := r.users[user.ID()]; ok {
		r.logger.Error("user already exists", "user_id", user.ID(), "email", user.Email())
		return models.UserID{}, errors.NewAlreadyExistsError("user with nickname %s already exists", user.Nickname()).
			WithDetails("nickname", user.Nickname())
	}

	r.users[user.ID()] = user
	r.logger.Info("user created", "user_id", user.ID(), "email", user.Email())
	return user.ID(), nil
}

func (r *UserRepository) Get(ctx context.Context, id models.UserID) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.logger.Debug("attempting to get user", "user_id", id)

	user, ok := r.users[id]
	if !ok {
		r.logger.Error("user not found", "user_id", id)
		return nil, errors.NewNotFoundError("user with id %s not found", id)
	}

	r.logger.Info("user found", "user_id", id, "email", user.Email())

	return user, nil
}

func (r *UserRepository) GetByNickname(ctx context.Context, nickname string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.logger.Debug("attempting to get user by nickname", "nickname", nickname)

	for _, user := range r.users {
		if user.Nickname() == nickname {
			r.logger.Info("user found", "user_id", user.ID(), "email", user.Email())
			return user, nil
		}
	}

	r.logger.Error("user not found", "nickname", nickname)
	return nil, errors.NewNotFoundError("user with nickname %s not found", nickname)
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debug("attempting to update user", "user_id", user.ID(), "email", user.Email())

	if _, ok := r.users[user.ID()]; !ok {
		r.logger.Error("user not found", "user_id", user.ID(), "email", user.Email())
		return errors.NewNotFoundError("user with id %s not found", user.ID())
	}

	r.users[user.ID()] = user
	r.logger.Info("user updated", "user_id", user.ID(), "email", user.Email())
	return nil
}
