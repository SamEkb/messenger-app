package friendship

import "context"

func (u *UseCase) AcceptFriendRequest(ctx context.Context, recipientID, requestorID string) error {
	u.logger.Info("accepting friend request")

	err := u.txManager.RunTx(ctx, func(txCtx context.Context) error {
		if err := u.friendRepository.AcceptFriendRequest(txCtx, recipientID, requestorID); err != nil {
			u.logger.Error("failed to accept friend request", "error", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	u.logger.Info("friend request accepted")

	return nil
}
