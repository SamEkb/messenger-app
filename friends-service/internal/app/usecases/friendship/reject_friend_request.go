package friendship

import "context"

func (u *UseCase) RejectFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	u.logger.Info("rejecting friend request")

	if err := u.friendRepository.RejectFriendRequest(ctx, recipientID, requestorID); err != nil {
		u.logger.Error("failed to reject friend request", "error", err)
		return err
	}

	u.logger.Info("friend request rejected")

	return nil
}
