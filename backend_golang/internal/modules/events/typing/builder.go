package typing

import (
	"time"
	"strings"

	"messenger/backend_golang/internal/common"
)

func BuildTypingEnvelope(chatType string,
	chatID string,
	targetID string,
	senderID string,
	username string,
) common.Envelope {

	return common.Envelope{
		Kind:      "event",
		Action:    "typing",
		ChatType:  strings.ToLower(strings.TrimSpace(chatType)),
		ChatID:    chatID,
		TargetID:  targetID,
		SenderID:  senderID,
		Timestamp: time.Now().UnixMilli(),
		Payload: map[string]any{
			"username": username,
		},
	}
}
