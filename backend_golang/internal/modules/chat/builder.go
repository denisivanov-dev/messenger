package chat

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/modules/utils"
)

func BuildMessage(env common.Envelope) common.Envelope {
	t := strings.TrimSpace(strings.ToLower(env.ChatType))
	chatID := utils.ResolveChatID(t, env.SenderID, env.TargetID)

	var payload common.MessagePayload
	b, _ := json.Marshal(env.Payload)
	_ = json.Unmarshal(b, &payload)

	if payload.MessageID == "" {
		payload.MessageID = uuid.NewString()
	}
	payload.EditedAt = 0

	return common.Envelope{
		Kind:      "message",
		Action:    "send",
		ChatID:    chatID,
		ChatType:  t,
		SenderID:  env.SenderID,
		TargetID:  env.TargetID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
}

func BuildSystemMessage(action, chatType, fromUserID, toUserID, text, context string) common.Envelope {
	t := strings.TrimSpace(strings.ToLower(chatType))
	chatID := utils.ResolveChatID(t, fromUserID, toUserID)

	sysPayload := common.SystemPayload{
		SystemID:  uuid.NewString(),
		Text:      text,
		Event:     action,
		ActorID:   fromUserID,
		TargetID:  toUserID,
		Timestamp: time.Now().UnixMilli(),
		Context:   context,
	}

	return common.Envelope{
		Kind:      "system",
		Action:    action,
		ChatID:    chatID,
		ChatType:  t,
		SenderID:  "0",
		TargetID:  toUserID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   sysPayload,
	}
}