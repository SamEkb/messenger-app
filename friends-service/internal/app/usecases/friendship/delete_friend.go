package friendship

import "context"

func (u *UseCase) DeleteFriend(ctx context.Context, userID, friendID string) error {
	u.logger.Info("deleting friend")

	err := u.txManager.RunTx(ctx, func(txCtx context.Context) error {
		if err := u.friendRepository.Delete(txCtx, userID, friendID); err != nil {
			u.logger.Error("failed to delete friend", "error", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	u.logger.Info("friend deleted")

	return nil
}
