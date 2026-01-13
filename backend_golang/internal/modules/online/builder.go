package online

import (
	"time"

	"messenger/backend_golang/internal/common"
)

func BuildStatusMessage(userID string, status string) common.Envelope {
	payload := common.OnlinePayload{
		Status: status, // "online" | "offline" | "away"
	}

	return common.Envelope{
		Kind:      "event",
		Action:    "user_status",
		SenderID:  userID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
}
