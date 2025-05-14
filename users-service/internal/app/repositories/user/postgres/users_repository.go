package postgres

import (
	"context"

	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/pkg/platform/postgres"
	"github.com/SamEkb/messenger-app/users-service/internal/app/models"
	"github.com/SamEkb/messenger-app/users-service/internal/app/ports"
)

var _ ports.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	txManager *postgres.TxManager
	logger    logger.Logger
}

func NewUserRepository(txManager *postgres.TxManager, logger logger.Logger) *UserRepository {
	return &UserRepository{
		txManager: txManager,
		logger:    logger.With("component", "user_repository"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (models.UserID, error) {
	r.logger.Debug("attempting to create new user", "user_id", user.ID(), "email", user.Email())

	q := r.txManager.GetQueryEngine(ctx)
	_, err := q.ExecContext(ctx, `
		INSERT INTO users (id, email, nickname, description, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID(), user.Email(), user.Nickname(), user.Description(), user.AvatarURL())
	if err != nil {
		r.logger.Error("failed to create user", "error", err)
		return models.UserID{}, errors.NewAlreadyExistsError("user with nickname %s or email %s already exists",
			user.Nickname(), user.Email())
	}

	r.logger.Info("user created", "user_id", user.ID(), "email", user.Email())
	return user.ID(), nil
}

func (r *UserRepository) Get(ctx context.Context, id models.UserID) (*models.User, error) {
	r.logger.Debug("attempting to get user", "user_id", id)

	q := r.txManager.GetQueryEngine(ctx)
	var user struct {
		ID          string `db:"id"`
		Email       string `db:"email"`
		Nickname    string `db:"nickname"`
		Description string `db:"description"`
		AvatarURL   string `db:"avatar_url"`
	}

	err := q.GetContext(ctx, &user, `
		SELECT id, email, nickname, description, avatar_url 
		FROM users 
		WHERE id = $1
	`, id)
	if err != nil {
		r.logger.Error("user not found", "user_id", id, "error", err)
		return nil, errors.NewNotFoundError("user with id %s not found", id)
	}

	userID, err := models.ParseUserID(user.ID)
	if err != nil {
		r.logger.Error("failed to parse user ID", "id", user.ID, "error", err)
		return nil, err
	}

	result, err := models.NewUser(userID, user.Email, user.Nickname, user.Description, user.AvatarURL)
	if err != nil {
		r.logger.Error("failed to create user model", "error", err)
		return nil, err
	}

	r.logger.Info("user found", "user_id", id, "email", user.Email)
	return result, nil
}

func (r *UserRepository) GetByNickname(ctx context.Context, nickname string) (*models.User, error) {
	r.logger.Debug("attempting to get user by nickname", "nickname", nickname)

	q := r.txManager.GetQueryEngine(ctx)
	var user struct {
		ID          string `db:"id"`
		Email       string `db:"email"`
		Nickname    string `db:"nickname"`
		Description string `db:"description"`
		AvatarURL   string `db:"avatar_url"`
	}

	err := q.GetContext(ctx, &user, `
		SELECT id, email, nickname, description, avatar_url 
		FROM users 
		WHERE nickname = $1
	`, nickname)
	if err != nil {
		r.logger.Error("user not found", "nickname", nickname, "error", err)
		return nil, errors.NewNotFoundError("user with nickname %s not found", nickname)
	}

	userID, err := models.ParseUserID(user.ID)
	if err != nil {
		r.logger.Error("failed to parse user ID", "id", user.ID, "error", err)
		return nil, err
	}

	result, err := models.NewUser(userID, user.Email, user.Nickname, user.Description, user.AvatarURL)
	if err != nil {
		r.logger.Error("failed to create user model", "error", err)
		return nil, err
	}

	r.logger.Info("user found", "nickname", nickname, "user_id", userID)
	return result, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	r.logger.Debug("attempting to update user", "user_id", user.ID(), "email", user.Email())

	q := r.txManager.GetQueryEngine(ctx)
	result, err := q.ExecContext(ctx, `
		UPDATE users 
		SET email = $1, nickname = $2, description = $3, avatar_url = $4 
		WHERE id = $5
	`, user.Email(), user.Nickname(), user.Description(), user.AvatarURL(), user.ID())
	if err != nil {
		r.logger.Error("failed to update user", "user_id", user.ID(), "error", err)
		return errors.NewInternalError(err, "failed to update user")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected", "error", err)
		return errors.NewInternalError(err, "failed to get rows affected")
	}

	if rows == 0 {
		r.logger.Error("user not found", "user_id", user.ID())
		return errors.NewNotFoundError("user with id %s not found", user.ID())
	}

	r.logger.Info("user updated", "user_id", user.ID(), "email", user.Email())
	return nil
}
