package postgres

import (
	"context"
	"time"

	"github.com/SamEkb/messenger-app/friends-service/config/env"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/models"
	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
	"github.com/SamEkb/messenger-app/pkg/platform/errors"
	"github.com/SamEkb/messenger-app/pkg/platform/logger"
	"github.com/SamEkb/messenger-app/pkg/platform/postgres"
	"github.com/google/uuid"
)

var _ ports.FriendshipRepository = (*FriendshipRepository)(nil)

type FriendshipRepository struct {
	txManager *postgres.TxManager
	db        *env.DBConfig
	logger    logger.Logger
}

func NewFriendshipRepository(txManager *postgres.TxManager, db *env.DBConfig, logger logger.Logger) *FriendshipRepository {
	return &FriendshipRepository{
		txManager: txManager,
		db:        db,
		logger:    logger.With("component", "friendship_repository"),
	}
}

func (r *FriendshipRepository) GetFriends(ctx context.Context, userID string) ([]*models.Friendship, error) {
	r.logger.Debug("getting friends", "user_id", userID)

	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > r.db.Timeout {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.db.Timeout)
		defer cancel()
	}

	q := r.txManager.GetQueryEngine(ctx)
	var friendships []struct {
		ID          string    `db:"id"`
		RequestorID string    `db:"requestor_id"`
		RecipientID string    `db:"recipient_id"`
		Status      string    `db:"status"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	err := q.GetContext(ctx, &friendships, `
		SELECT id, requestor_id, recipient_id, status, created_at, updated_at 
		FROM friendships 
		WHERE (requestor_id = $1 OR recipient_id = $1) AND status = $2
	`, userID, models.FriendshipStatusAccepted)

	if err != nil {
		r.logger.Error("failed to get friends", "error", err, "user_id", userID)
		return nil, errors.NewInternalError(err, "failed to get friends")
	}

	result := make([]*models.Friendship, 0, len(friendships))
	for _, f := range friendships {
		friendship, err := r.mapToModel(f.ID, f.RequestorID, f.RecipientID, f.Status, f.CreatedAt, f.UpdatedAt)
		if err != nil {
			r.logger.Error("failed to map friendship", "error", err)
			continue
		}
		result = append(result, friendship)
	}

	r.logger.Debug("got friends", "user_id", userID, "count", len(result))
	return result, nil
}

func (r *FriendshipRepository) SendFriendRequest(ctx context.Context, requestorID, recipientID string) error {
	r.logger.Debug("sending friend request", "requestor_id", requestorID, "recipient_id", recipientID)

	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > r.db.Timeout {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.db.Timeout)
		defer cancel()
	}

	if requestorID == recipientID {
		return errors.NewInvalidInputError("cannot send friend request to yourself")
	}

	friendship, err := models.NewFriendship(requestorID, recipientID)
	if err != nil {
		r.logger.Error("failed to create friendship model", "error", err)
		return err
	}

	q := r.txManager.GetQueryEngine(ctx)
	_, err = q.ExecContext(ctx, `
		INSERT INTO friendships (id, requestor_id, recipient_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, friendship.ID(), friendship.RequestorID(), friendship.RecipientID(),
		friendship.Status(), friendship.CreatedAt(), friendship.UpdatedAt())

	if err != nil {
		r.logger.Error("failed to send friend request", "error", err)
		return errors.NewAlreadyExistsError("friend request already exists")
	}

	r.logger.Info("friend request sent", "requestor_id", requestorID, "recipient_id", recipientID)
	return nil
}

func (r *FriendshipRepository) AcceptFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	r.logger.Debug("accepting friend request", "recipient_id", recipientID, "requestor_id", requestorID)

	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > r.db.Timeout {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.db.Timeout)
		defer cancel()
	}

	q := r.txManager.GetQueryEngine(ctx)
	var friendship struct {
		ID          string    `db:"id"`
		RequestorID string    `db:"requestor_id"`
		RecipientID string    `db:"recipient_id"`
		Status      string    `db:"status"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	err := q.GetContext(ctx, &friendship, `
		SELECT id, requestor_id, recipient_id, status, created_at, updated_at 
		FROM friendships 
		WHERE requestor_id = $1 AND recipient_id = $2 AND status = $3
	`, requestorID, recipientID, models.FriendshipStatusRequested)

	if err != nil {
		r.logger.Error("friend request not found", "error", err)
		return errors.NewNotFoundError("friend request not found")
	}

	friendshipModel, err := r.mapToModel(
		friendship.ID,
		friendship.RequestorID,
		friendship.RecipientID,
		friendship.Status,
		friendship.CreatedAt,
		friendship.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("failed to map friendship", "error", err)
		return err
	}

	friendshipModel.Accept()

	_, err = q.ExecContext(ctx, `
		UPDATE friendships 
		SET status = $1, updated_at = $2 
		WHERE id = $3
	`, friendshipModel.Status(), friendshipModel.UpdatedAt(), friendshipModel.ID())

	if err != nil {
		r.logger.Error("failed to accept friend request", "error", err)
		return errors.NewInternalError(err, "failed to accept friend request")
	}

	r.logger.Info("friend request accepted", "recipient_id", recipientID, "requestor_id", requestorID)
	return nil
}

func (r *FriendshipRepository) RejectFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	r.logger.Debug("rejecting friend request", "recipient_id", recipientID, "requestor_id", requestorID)

	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > r.db.Timeout {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.db.Timeout)
		defer cancel()
	}

	q := r.txManager.GetQueryEngine(ctx)
	var friendship struct {
		ID          string    `db:"id"`
		RequestorID string    `db:"requestor_id"`
		RecipientID string    `db:"recipient_id"`
		Status      string    `db:"status"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	err := q.GetContext(ctx, &friendship, `
		SELECT id, requestor_id, recipient_id, status, created_at, updated_at 
		FROM friendships 
		WHERE requestor_id = $1 AND recipient_id = $2 AND status = $3
	`, requestorID, recipientID, models.FriendshipStatusRequested)

	if err != nil {
		r.logger.Error("friend request not found", "error", err)
		return errors.NewNotFoundError("friend request not found")
	}

	friendshipModel, err := r.mapToModel(
		friendship.ID,
		friendship.RequestorID,
		friendship.RecipientID,
		friendship.Status,
		friendship.CreatedAt,
		friendship.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("failed to map friendship", "error", err)
		return err
	}

	friendshipModel.Reject()

	_, err = q.ExecContext(ctx, `
		UPDATE friendships 
		SET status = $1, updated_at = $2 
		WHERE id = $3
	`, friendshipModel.Status(), friendshipModel.UpdatedAt(), friendshipModel.ID())

	if err != nil {
		r.logger.Error("failed to reject friend request", "error", err)
		return errors.NewInternalError(err, "failed to reject friend request")
	}

	r.logger.Info("friend request rejected", "recipient_id", recipientID, "requestor_id", requestorID)
	return nil
}

func (r *FriendshipRepository) Delete(ctx context.Context, userID string, friendID string) error {
	r.logger.Debug("deleting friendship", "user_id", userID, "friend_id", friendID)

	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > r.db.Timeout {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.db.Timeout)
		defer cancel()
	}

	q := r.txManager.GetQueryEngine(ctx)
	result, err := q.ExecContext(ctx, `
		DELETE FROM friendships 
		WHERE ((requestor_id = $1 AND recipient_id = $2) OR (requestor_id = $2 AND recipient_id = $1)) 
		AND status = $3
	`, userID, friendID, models.FriendshipStatusAccepted)

	if err != nil {
		r.logger.Error("failed to delete friendship", "error", err)
		return errors.NewInternalError(err, "failed to delete friendship")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected", "error", err)
		return errors.NewInternalError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		r.logger.Error("friendship not found")
		return errors.NewNotFoundError("friendship not found")
	}

	r.logger.Info("friendship deleted", "user_id", userID, "friend_id", friendID)
	return nil
}

func (r *FriendshipRepository) mapToModel(id, requestorID, recipientID, status string, createdAt, updatedAt time.Time) (*models.Friendship, error) {
	friendshipID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return models.NewFriendshipFromDB(
		friendshipID,
		requestorID,
		recipientID,
		models.FriendshipStatus(status),
		createdAt,
		updatedAt,
	), nil
}
