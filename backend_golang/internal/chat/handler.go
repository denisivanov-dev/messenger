package chat

import (
	"log"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
)

func HandleSendMessage(in common.IncomingSendMessage, userID, username string, rdb *redis.Client) (common.OutgoingMessage, bool) {
	outMsg := BuildMessage(in, userID, username)

	if outMsg.ChatID == "" {
		log.Printf("empty ChatID in message, skipping: %+v", outMsg)
		return common.OutgoingMessage{}, false
	}

	go SaveMessageToRedisHistory(rdb, outMsg.ChatID, outMsg)

	return outMsg, true
}