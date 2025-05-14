package in_memory

import (
	"context"
	"sync"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
)

var _ ports.AuthRepository = (*AuthRepository)(nil)

type AuthRepository struct {
	mx      sync.Mutex
	storage map[models.UserID]*models.User
	logger  logger.Logger
}

func NewAuthRepository(logger logger.Logger) *AuthRepository {
	return &AuthRepository{
		storage: make(map[models.UserID]*models.User),
		logger:  logger.With("component", "auth_repository"),
	}
}

func (r *AuthRepository) Create(ctx context.Context, user *models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Debug("attempting to create new user", "user_id", user.ID(), "email", user.Email())

	existingUser, _ := r.FindUserByEmail(ctx, user.Email())
	if existingUser != nil {
		r.logger.Warn("user with this email already exists", "email", user.Email())
		return errors.NewAlreadyExistsError("user with email %s already exists", user.Email()).
			WithDetails("email", user.Email())
	}

	r.storage[user.ID()] = user
	r.logger.Info("user created successfully", "user_id", user.ID(), "email", user.Email())
	return nil
}

func (r *AuthRepository) FindUserByID(ctx context.Context, userID models.UserID) (*models.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Debug("looking for user by ID", "user_id", userID)

	user, ok := r.storage[userID]
	if !ok {
		r.logger.Debug("user not found", "user_id", userID)
		return nil, errors.NewNotFoundError("user with ID %s not found", userID.String()).
			WithDetails("user_id", userID.String())
	}

	r.logger.Debug("user found", "user_id", userID, "email", user.Email())
	return user, nil
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Debug("looking for user by email", "email", email)

	var user *models.User
	for id := range r.storage {
		currentUser := r.storage[id]
		if r.storage[id].Email() == email {
			user = currentUser
			break
		}
	}

	if user == nil {
		r.logger.Debug("user not found", "email", email)
		return nil, errors.NewNotFoundError("user with email %s not found", email).
			WithDetails("email", email)
	}

	r.logger.Debug("user found", "user_id", user.ID(), "email", email)
	return user, nil
}

func (r *AuthRepository) Update(ctx context.Context, user *models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.logger.Debug("attempting to update user", "user_id", user.ID())

	if _, ok := r.storage[user.ID()]; !ok {
		r.logger.Debug("user not found for update", "user_id", user.ID())
		return errors.NewNotFoundError("user with ID %s not found", user.ID().String()).
			WithDetails("user_id", user.ID().String())
	}

	r.storage[user.ID()] = user
	r.logger.Info("user updated successfully", "user_id", user.ID())

	return nil
}
