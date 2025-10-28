package chat

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/modules/voice"
)

func HandleSendMessage(env common.Envelope, rdb *redis.Client) (common.Envelope, bool) {
	outMsg := BuildMessage(env) 

	if outMsg.ChatID == "" {
		log.Printf("[chat] empty ChatID in message, skipping: %+v", outMsg)
		return common.Envelope{}, false
	}

	go SaveMessageToRedisHistory(rdb, outMsg.ChatID, outMsg)

	return outMsg, true
}

func HandleSystemMessage(action, chatType, fromUserID, toUserID, text, context string, rdb *redis.Client) (common.Envelope, bool) {
	outMsg := BuildSystemMessage(action, chatType, fromUserID, toUserID, text, context)

	if outMsg.ChatID == "" {
		log.Printf("[chat] empty ChatID in system message, skipping: %+v", outMsg)
		return common.Envelope{}, false
	}

	if action == "call_started" {
		roomID := outMsg.ChatID
		payload, ok := outMsg.Payload.(common.SystemPayload)
		if !ok {
			log.Printf("[chat] invalid system payload for call_started")
			return outMsg, false
		}

		voice.SetOngoingCallMessageID(rdb, roomID, payload.SystemID)
		voice.SetCallStartTime(rdb, roomID, time.Now().Unix())
	}

	go SaveMessageToRedisHistory(rdb, outMsg.ChatID, outMsg)
	return outMsg, true
}
