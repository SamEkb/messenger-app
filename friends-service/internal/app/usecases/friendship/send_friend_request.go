package friendship

import "context"

func (u *UseCase) SendFriendRequest(ctx context.Context, requestorID, recipientID string) error {
	u.logger.Info("sending friend request")

	if err := u.friendRepository.SendFriendRequest(ctx, requestorID, recipientID); err != nil {
		u.logger.Error("failed to send friend request", "error", err)
		return err
	}

	u.logger.Info("friend request sent")

	return nil
}
