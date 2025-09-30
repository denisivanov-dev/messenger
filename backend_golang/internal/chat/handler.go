package chat

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/voice"
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

func HandleSystemMessage(msgType string, chatType string, fromUserID string, toUserID string, callInfo *common.CallInfo, rdb *redis.Client) (common.OutgoingMessage, bool) {
	outMsg := BuildSystemMessage(msgType, chatType, fromUserID, toUserID, callInfo)

	if outMsg.ChatID == "" {
		log.Printf("empty ChatID in system message, skipping: %+v", outMsg)
		return common.OutgoingMessage{}, false
	}

	if msgType == "call_started" && callInfo != nil {
		roomID := outMsg.ChatID

		voice.SetOngoingCallMessageID(rdb, roomID, outMsg.MessageID)

		if callInfo.StartedAt > 0 {
			voice.SetCallStartTime(rdb, roomID, callInfo.StartedAt)
		} else {
			voice.SetCallStartTime(rdb, roomID, time.Now().Unix())
		}
	}

	go SaveMessageToRedisHistory(rdb, outMsg.ChatID, outMsg)

	return outMsg, true
}