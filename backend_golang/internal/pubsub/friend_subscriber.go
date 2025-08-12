package pubsub

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/ws"
)

func StartFriendRequestSubscriber(ctx context.Context, rdb *redis.Client, hub *ws.Hub) {
	pubsub := rdb.Subscribe(ctx, "friend_requests")
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var payload common.FriendRequestPayload
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Printf("Failed to parse friend request: %v", err)
				continue
			}

			if payload.Type == "friend_request_sent" {
				log.Printf("Sent friend request from %s to %s", payload.FromID, payload.ToID)

				hub.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "outgoing", // для отправителя
				})

				hub.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "incoming", // для получателя
				})
			}

			if payload.Type == "friend_request_canceled" {
				log.Printf("Canceled friend request from %s to %s", payload.FromID, payload.ToID)

				hub.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "none", // отправителю: "ничего нет"
				})

				hub.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "none", // получателю: "ничего нет"
				})
			}

			if payload.Type == "friend_request_accepted" {
				log.Printf("Friend request accepted: %s ↔ %s", payload.FromID, payload.ToID)

				hub.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "friends", // для отправителя
				})

				hub.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "friends", // для получателя
				})
			}

			if payload.Type == "friend_request_declined" {
				log.Printf("Friend request declined: %s ❌ %s", payload.FromID, payload.ToID)

				hub.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "none", // отправитель больше не видит
				})

				hub.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "none", // получатель тоже ничего не видит
				})
			}

			if payload.Type == "friend_removed" {
				log.Printf("Friend removed: %s ❌ %s", payload.FromID, payload.ToID)

				hub.SendToUser(payload.FromID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.ToID,
					"status":  "none", // Удалившему: больше не друг
				})

				hub.SendToUser(payload.ToID, map[string]any{
					"type":    "friend_request_update",
					"user_id": payload.FromID,
					"status":  "none", // Удалённому: тоже больше не друг
				})
			}
		}
	}()
}