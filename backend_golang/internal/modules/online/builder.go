package online

import (
	"encoding/json"
	"time"

	"messenger/backend_golang/internal/common"
)

func BuildStatusMessage(userID string, status string) []byte {
	payload := common.OnlinePayload{
		UserID: userID,
		Status: status, // "online" | "offline" | "away"
	}

	env := common.Envelope{
		Kind:      "event",
		Action:    "user_status",
		SenderID:  userID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}

	data, _ := json.Marshal(env)
	return data
}
