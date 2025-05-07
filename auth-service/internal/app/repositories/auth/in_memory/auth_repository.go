package in_memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
)

var _ ports.AuthRepository = (*AuthRepository)(nil)

type AuthRepository struct {
	mx      sync.Mutex
	storage map[models.UserID]*models.User
}

var (
	_ ports.AuthRepository = (*AuthRepository)(nil)
)

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		storage: make(map[models.UserID]*models.User),
	}
}

func (r *AuthRepository) Create(ctx context.Context, user *models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if user, _ := r.FindUserByEmail(ctx, user.Email()); user != nil {
		return errors.New("user with this username already exists")
	}

	r.storage[user.ID()] = user
	return nil
}

func (r *AuthRepository) FindUserByID(ctx context.Context, userID models.UserID) (*models.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	user, ok := r.storage[userID]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, username string) (*models.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	var user *models.User
	for id := range r.storage {
		currentUser := r.storage[id]
		if r.storage[id].Username() == username {
			user = currentUser
			break
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user with this username: %s not found", username)
	}

	return user, nil
}

func (r *AuthRepository) Update(ctx context.Context, user *models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.storage[user.ID()]; !ok {
		return errors.New("user not found")
	}

	r.storage[user.ID()] = user

	return nil
}
