package chat

import (
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
)

func HandleSendMessage(in common.IncomingSendMessage, userID, username string, rdb *redis.Client) ([]byte, bool) {
	outMsg := BuildMessage(in, userID, username)

	if outMsg.ChatID == "" {
		log.Printf("empty ChatID in message, skipping: %+v", outMsg)
		return nil, false
	}

	go SaveMessageToRedisHistory(rdb, outMsg.ChatID, outMsg)

	outBytes, err := json.Marshal(outMsg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return nil, false
	}

	return outBytes, true
}