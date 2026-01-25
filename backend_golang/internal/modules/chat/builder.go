package chat

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"messenger/backend_golang/internal/common"
)

func BuildMessageEnvelope(
	action string,
	chatType string,
	chatID string,
	senderID string,
	targetID string,
	username string,
	payload common.MessagePayload,
) common.Envelope {

	t := strings.ToLower(strings.TrimSpace(chatType))

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
		ChatType:  t,
		ChatID:    chatID,
		SenderID:  senderID,
		TargetID:  targetID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
}

func BuildSystemMessageEnvelope(
	action string,
	chatType string,
	chatID string,
	fromUserID string,
	toUserID string,
	text string,
	context string,
) common.Envelope {

	t := strings.ToLower(strings.TrimSpace(chatType))

	payload := common.SystemPayload{
		SystemID:  uuid.NewString(),
		Event:     action,
		Text:      text,
		ActorID:   fromUserID,
		TargetID:  toUserID,
		Context:   context,
		Timestamp: time.Now().UnixMilli(),
	}

	return common.Envelope{
		Kind:      "system",
		Action:    action,
		ChatType:  t,
		ChatID:    chatID,
		SenderID:  "0",
		TargetID:  toUserID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
}
