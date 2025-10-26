package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	rds "github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/modules/utils"
)

func GetRoomKey(userID, chatType, receiverID string) string {
	if chatType == "private" {
		return utils.GeneratePrivateChatKey(userID, receiverID)
	}
	return GlobalRoomID 
}

// ResolveRoom checks if the user has access to the chat.
func ResolveRoom(rdb *rds.Client, userID, chatType, receiverID string) (string, bool) {
	ctx := context.Background()

	chatKey := GetRoomKey(userID, chatType, receiverID)

	if chatKey == GlobalRoomID {
		return chatKey, true
	}
	
	chatID, err := rdb.Get(ctx, fmt.Sprintf("chat_id:%s", chatKey)).Result()
	if err != nil {
		log.Printf("[ResolveRoom] Redis GET chat_id:%s failed: %v", chatKey, err)
		return "", false
	}

	log.Printf("[ResolveRoom] Mapped chatKey %s → chatID %s", chatKey, chatID)

	ok, err := rdb.SIsMember(ctx, fmt.Sprintf("chat:%s:participants", chatID), userID).Result()
	if err != nil {
		log.Printf("[ResolveRoom] Redis SISMEMBER failed for chat:%s: %v", chatID, err)
		return "", false
	}
	if !ok {
		log.Printf("[ResolveRoom] Access denied: user %s not in chat %s", userID, chatID)
		return "", false
	}

	return chatID, true
}

func SendHistory(rdb *rds.Client, roomID string, send chan []byte, limit int) {
	msgs, err := LoadMessageHistory(rdb, roomID, int64(limit))
	if err != nil {
		log.Printf("failed to load chat history for %s: %v", roomID, err)
		return
	}

	for _, msg := range msgs {
		if out, err := json.Marshal(msg); err == nil {
			send <- out
		}
	}
}