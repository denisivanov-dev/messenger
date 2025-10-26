package utils

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/gateway/types"
)

// unmarshalPayload safely unmarshals incoming JSON payloads into a Go struct.
func UnmarshalPayload(raw []byte, dest any, c types.ClientLike, msg string) bool {
	if err := json.Unmarshal(raw, dest); err != nil {
		c.SendError(msg)
		log.Printf("[unmarshalPayload] %s: %v", msg, err)

		return false
	}
	return true
}
