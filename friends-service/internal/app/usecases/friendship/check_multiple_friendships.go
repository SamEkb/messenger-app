package friendship

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
)

func (u *UseCase) CheckMultipleFriendships(ctx context.Context, userIDs []string) ([]ports.UserPair, error) {
	u.logger.Info("checking multiple friendships", "user_count", len(userIDs))

	nonFriendPairs := make([]ports.UserPair, 0)

	for i := 0; i < len(userIDs); i++ {
		for j := i + 1; j < len(userIDs); j++ {
			user1 := userIDs[i]
			user2 := userIDs[j]

			isFriend, err := u.areFriends(ctx, user1, user2)
			if err != nil {
				u.logger.Error("failed to check friendship status", "user1", user1, "user2", user2, "error", err)
				return nil, err
			}

			if !isFriend {
				nonFriendPairs = append(nonFriendPairs, ports.UserPair{
					UserID1: user1,
					UserID2: user2,
				})
			}
		}
	}

	u.logger.Info("multiple friendships check completed",
		"total_pairs", (len(userIDs)*(len(userIDs)-1))/2,
		"non_friend_pairs", len(nonFriendPairs))

	return nonFriendPairs, nil
}

func (u *UseCase) areFriends(ctx context.Context, user1ID, user2ID string) (bool, error) {
	friends, err := u.friendRepository.GetFriends(ctx, user1ID)
	if err != nil {
		return false, err
	}

	for _, friendship := range friends {
		if friendship.IsAccepted() {
			if (friendship.RequestorID() == user1ID && friendship.RecipientID() == user2ID) ||
				(friendship.RequestorID() == user2ID && friendship.RecipientID() == user1ID) {
				return true, nil
			}
		}
	}

	return false, nil
}
