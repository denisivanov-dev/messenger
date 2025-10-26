package online

import (
	"encoding/json"

	"messenger/backend_golang/internal/common"
)
func BuildStatusMessage(userID string, status common.Status) []byte {
	msg := common.StatusMessage{
		Event:  "user_status",
		UserID: userID,
		Status: status,
	}
	b, _ := json.Marshal(msg)

	return b
}
