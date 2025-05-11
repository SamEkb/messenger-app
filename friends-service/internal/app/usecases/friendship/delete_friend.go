package friendship

import "context"

func (u *UseCase) DeleteFriend(ctx context.Context, userID, friendID string) error {
	u.logger.Info("deleting friend")

	if err := u.friendRepository.Delete(ctx, userID, friendID); err != nil {
		u.logger.Error("failed to delete friend", "error", err)
		return err
	}

	u.logger.Info("friend deleted")

	return nil
}
