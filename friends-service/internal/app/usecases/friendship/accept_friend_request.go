package friendship

import "context"

func (u *UseCase) AcceptFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	u.logger.Info("accepting friend request")

	if err := u.friendRepository.AcceptFriendRequest(ctx, recipientID, requestorID); err != nil {
		u.logger.Error("failed to accept friend request", "error", err)
		return err
	}

	u.logger.Info("friend request accepted")

	return nil
}
