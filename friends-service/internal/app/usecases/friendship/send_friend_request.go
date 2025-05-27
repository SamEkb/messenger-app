package friendship

import "context"

func (u *UseCase) SendFriendRequest(ctx context.Context, requestorID, recipientID string) error {
	u.logger.Info("sending friend request")

	err := u.txManager.RunTx(ctx, func(txCtx context.Context) error {
		if err := u.friendRepository.SendFriendRequest(txCtx, requestorID, recipientID); err != nil {
			u.logger.Error("failed to send friend request", "error", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	u.logger.Info("friend request sent")

	return nil
}
