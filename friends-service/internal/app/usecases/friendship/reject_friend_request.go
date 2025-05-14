package friendship

import "context"

func (u *UseCase) RejectFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	u.logger.Info("rejecting friend request")

	err := u.txManager.RunTx(ctx, func(txCtx context.Context) error {
		if err := u.friendRepository.RejectFriendRequest(txCtx, recipientID, requestorID); err != nil {
			u.logger.Error("failed to reject friend request", "error", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	u.logger.Info("friend request rejected")

	return nil
}
