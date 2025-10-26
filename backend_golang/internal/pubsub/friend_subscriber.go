package pubsub

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/types"
)

// optional interface for hubs that support direct user messaging
type userSender interface {
	SendToUser(userID string, data any)
}

func StartFriendRequestSubscriber(ctx context.Context, rdb *redis.Client, hub types.HubLike) {
	pubsub := rdb.Subscribe(ctx, "friend_requests")
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var payload common.FriendRequestPayload
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Printf("Failed to parse friend request: %v", err)
				continue
			}

			sender, ok := hub.(userSender)
			if !ok {
				log.Printf("Hub does not implement SendToUser()")
				continue
			}

			if payload.Type == "friend_request_sent" {
				log.Printf("Sent friend request from %s to %s", payload.FromID, payload.ToID)

				sender.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "outgoing", // for sender
				})

				sender.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "incoming", // for receiver
				})
			}

			if payload.Type == "friend_request_canceled" {
				log.Printf("Canceled friend request from %s to %s", payload.FromID, payload.ToID)

				sender.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "none", // for sender: no active request
				})

				sender.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "none", // for receiver: no active request
				})
			}

			if payload.Type == "friend_request_accepted" {
				log.Printf("Friend request accepted: %s <-> %s", payload.FromID, payload.ToID)

				sender.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "friends", // for sender
				})

				sender.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "friends", // for receiver
				})
			}

			if payload.Type == "friend_request_declined" {
				log.Printf("Friend request declined: %s %s", payload.FromID, payload.ToID)

				sender.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "none", // sender no longer sees the request
				})

				sender.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "none", // receiver no longer sees the request
				})
			}

			if payload.Type == "friend_removed" {
				log.Printf("Friend removed: %s %s", payload.FromID, payload.ToID)

				sender.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "none", // for the user who removed the friend
				})

				sender.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "none", // for the removed user
				})
			}
		}
	}()
}