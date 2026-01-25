package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	rds "github.com/redis/go-redis/v9"
	
	"messenger/backend_golang/internal/modules/utils"
	"messenger/backend_golang/internal/modules/constants"
)

// ResolveRoom checks if the user has access to the chat.
func ResolveRoom(rdb *rds.Client, userID, chatType, receiverID string) (string, bool) {
	ctx := context.Background()

	chatID := utils.ResolveChatID(chatType, userID, receiverID)
	if chatID == "" {
		log.Printf("[ResolveRoom] empty chatID (chatType=%s user=%s receiver=%s)", chatType, userID, receiverID)
		return "", false
	}

	if chatID == constants.GlobalRoomID {
		return chatID, true
	}

	ok, err := rdb.SIsMember(ctx, fmt.Sprintf("chat:%s:participants", chatID), userID).Result()
	if err != nil {
		log.Printf("[ResolveRoom] Redis SISMEMBER failed for chat:%s: %v", chatID, err)
		return "", false
	}
	if !ok {
		log.Printf("[ResolveRoom] access denied: user %s not in chat %s", userID, chatID)
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
