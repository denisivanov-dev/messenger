package chat

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/modules/utils"
)

func BuildMessage(action, senderID, username string, env common.Envelope) common.Envelope {
	t := strings.TrimSpace(strings.ToLower(env.ChatType))
	chatID := utils.ResolveChatID(t, senderID, env.TargetID)

	payload := common.MessagePayload{
		Attachments: []common.Attachment{},
		EditedAt:    0,
		IsEdited:    false,
		MessageID:   uuid.NewString(),
		Pinned:      false,
		ReplyTo:     nil,
		ReplyToText: nil,
		ReplyToUser: nil,
		Text:        "",
		Username:    username,
	}

	raw, _ := json.Marshal(env.Payload)
	_ = json.Unmarshal(raw, &payload)

	if payload.MessageID == "" {
		payload.MessageID = uuid.NewString()
	}
	if payload.Username == "" {
		payload.Username = username
	}

	switch action {
	case "send":
		payload.IsEdited = false
		payload.EditedAt = 0

	case "edit":
		payload.IsEdited = true
		payload.EditedAt = time.Now().UnixMilli()

	case "delete":
		payload = common.MessagePayload{
			MessageID: payload.MessageID,
			Username:  username,
		}

	case "pin":
		payload.Pinned = true

	case "unpin":
		payload.Pinned = false

	case "reply":
		if payload.ReplyToUser == nil || *payload.ReplyToUser == "" {
			u := senderID
			payload.ReplyToUser = &u
		}
	}

	return common.Envelope{
		Kind:      "message",
		Action:    action,
		ChatID:    chatID,
		ChatType:  t,
		SenderID:  senderID,
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