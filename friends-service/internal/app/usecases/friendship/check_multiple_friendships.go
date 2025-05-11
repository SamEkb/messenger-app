package friendship

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/internal/app/ports"
)

func (u *UseCase) CheckMultipleFriendships(ctx context.Context, userIDs []string) ([]ports.UserPair, error) {
	friendshipMap := make(map[string]map[string]bool)

	for _, userID := range userIDs {
		friends, err := u.friendRepository.GetFriends(ctx, userID)
		if err != nil {
			return nil, err
		}

		if _, ok := friendshipMap[userID]; !ok {
			friendshipMap[userID] = make(map[string]bool)
		}

		for _, friend := range friends {
			if friend.IsAccepted() {
				otherID := friend.RecipientID()
				if otherID == userID {
					otherID = friend.RequestorID()
				}
				friendshipMap[userID][otherID] = true
			}
		}
	}

	var nonFriendPairs []ports.UserPair
	for i := 0; i < len(userIDs); i++ {
		for j := i + 1; j < len(userIDs); j++ {
			user1 := userIDs[i]
			user2 := userIDs[j]

			if friends, ok := friendshipMap[user1]; ok {
				if _, areFriends := friends[user2]; areFriends {
					continue
				}
			}

			nonFriendPairs = append(nonFriendPairs, ports.UserPair{
				UserID1: user1,
				UserID2: user2,
			})
		}
	}

	return nonFriendPairs, nil
}
