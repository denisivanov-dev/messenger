package chat

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/modules/utils"
	"messenger/backend_golang/internal/modules/voice"
)

func HandleSendUserMessage(env common.Envelope, chatID string, rdb *redis.Client,
	senderID, username string,) (common.Envelope, bool) {

	if chatID == "" {
		log.Printf("[chat] empty chatID in HandleSendUserMessage")
		return common.Envelope{}, false
	}

	payload, ok := env.Payload.(common.MessagePayload)
	if !ok {
		log.Printf("[chat] invalid message payload type")
		return common.Envelope{}, false
	}

	outMsg := BuildMessageEnvelope(
		env.Action,
		env.ChatType,
		chatID,
		senderID,
		env.TargetID,
		username,
		payload,
	)

	go SaveMessageToRedisHistory(rdb, chatID, outMsg)
	return outMsg, true
}

func HandleSystemMessage(
	action string,
	chatType string,
	fromUserID string,
	toUserID string,
	text string,
	context string,
	rdb *redis.Client,
) (common.Envelope, bool) {

	chatID := utils.ResolveChatID(chatType, fromUserID, toUserID)
	if chatID == "" {
		log.Printf("[chat] empty ChatID in system message, skipping")
		return common.Envelope{}, false
	}

	outMsg := BuildSystemMessageEnvelope(
		action,
		chatType,
		chatID,
		fromUserID,
		toUserID,
		text,
		context,
	)

	if action == "call_started" {
		payload, ok := outMsg.Payload.(common.SystemPayload)
		if !ok {
			log.Printf("[chat] invalid system payload for call_started")
			return outMsg, false
		}

		voice.SetOngoingCallMessageID(rdb, chatID, payload.SystemID)
		voice.SetCallStartTime(rdb, chatID, time.Now().Unix())
	}

	go SaveMessageToRedisHistory(rdb, outMsg.ChatID, outMsg)

	return outMsg, true
}