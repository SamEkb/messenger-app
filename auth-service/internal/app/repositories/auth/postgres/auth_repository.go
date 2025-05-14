package postgres

import (
	"context"

	"github.com/SamEkb/messenger-app/auth-service/internal/app/models"
	"github.com/SamEkb/messenger-app/auth-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/pkg/platform/postgres"
)

var _ ports.AuthRepository = (*AuthRepository)(nil)

type AuthRepository struct {
	q      postgres.QueryEngine
	logger logger.Logger
}

func NewAuthRepository(q postgres.QueryEngine, logger logger.Logger) *AuthRepository {
	return &AuthRepository{
		q:      q,
		logger: logger.With("component", "auth_repository"),
	}
}

func (r *AuthRepository) Create(ctx context.Context, user *models.User) error {
	r.logger.Debug("attempting to create new user", "user_id", user.ID(), "email", user.Email())

	_, err := r.q.ExecContext(ctx, `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
	`, user.ID(), user.Username(), user.Email(), user.Password())
	if err != nil {
		r.logger.Warn("user with this email already exists", "email", user.Email())
		return errors.NewAlreadyExistsError("user with email %s already exists", user.Email())
	}

	r.logger.Debug("user created successfully", "user_id", user.ID(), "email", user.Email())

	return nil
}

func (r *AuthRepository) FindUserByID(ctx context.Context, userID models.UserID) (*models.User, error) {
	r.logger.Debug("looking for user by ID", "user_id", userID)
	var user struct {
		ID       string `db:"id"`
		Username string `db:"username"`
		Email    string `db:"email"`
		Password string `db:"password"`
	}
	err := r.q.GetContext(ctx, &user, `
		SELECT id, username, email, password FROM users WHERE id = $1
	`, userID)
	if err != nil {
		r.logger.Warn("user not found", "user_id", userID)
		return nil, errors.NewNotFoundError("user with ID %s not found", userID.String())
	}

	uid, err := models.UserIDFromString(user.ID)
	if err != nil {
		r.logger.Warn("failed to parse UUID", "UUID", userID)
		return nil, err
	}

	r.logger.Debug("user found", "user_id", userID, "email", user.Email)

	return models.NewUser(uid, user.Username, user.Email, []byte(user.Password))
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.logger.Debug("looking for user by email", "email", email)

	var user struct {
		ID       string `db:"id"`
		Username string `db:"username"`
		Email    string `db:"email"`
		Password string `db:"password"`
	}
	err := r.q.GetContext(ctx, &user, `
		SELECT id, username, email, password FROM users WHERE email = $1
	`, email)
	if err != nil {
		r.logger.Warn("user not found", "email", email)
		return nil, errors.NewNotFoundError("user with email %s not found", email)
	}
	uid, err := models.UserIDFromString(user.ID)
	if err != nil {
		r.logger.Warn("failed to parse UUID", "UUID", user.ID)
		return nil, err
	}

	r.logger.Debug("user found", "user_id", user.ID, "email", email)
	return models.NewUser(uid, user.Username, user.Email, []byte(user.Password))
}

func (r *AuthRepository) Update(ctx context.Context, user *models.User) error {
	r.logger.Debug("attempting to update user", "user_id", user.ID())

	_, err := r.q.ExecContext(ctx, `
		UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4
	`, user.Username(), user.Email(), user.Password(), user.ID())
	if err != nil {
		r.logger.Warn("user not found for update", "user_id", user.ID())
		return errors.NewNotFoundError("user with ID %s not found", user.ID().String())
	}

	return nil
}
