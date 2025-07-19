package chat

import (
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

func HandleIncoming(raw []byte, userID, username string, rdb *redis.Client) (string, []byte, bool) {
	incMsg, ok := Parse(raw)
	if !ok {
		return "", nil, false
	}

	outMsg := BuildMessage(incMsg, userID, username)

	if outMsg.ChatID == "" {
		log.Printf("empty ChatID in message, skipping: %+v", outMsg)
		return "", nil, false
	}

	roomID := outMsg.ChatID
	go SaveMessageToRedisHistory(rdb, roomID, outMsg)

	outBytes, err := json.Marshal(outMsg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return "", nil, false
	}

	return roomID, outBytes, true
}
