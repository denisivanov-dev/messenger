package utils

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/types"
)

func UnmarshalPayload(raw []byte, dest any, c types.ClientLike, msg string) bool {
	var env common.Envelope
	if err := json.Unmarshal(raw, &env); err == nil && env.Payload != nil {
		data, err := json.Marshal(env.Payload)
		if err != nil {
			c.SendError(msg)
			log.Printf("[UnmarshalPayload] failed to marshal payload: %v", err)
			return false
		}
		if err := json.Unmarshal(data, dest); err != nil {
			c.SendError(msg)
			log.Printf("[UnmarshalPayload] failed to unmarshal payload: %v", err)
			return false
		}
		return true
	}

	// If Envelope doesnt match - try legacy version
	if err := json.Unmarshal(raw, dest); err != nil {
		c.SendError(msg)
		log.Printf("[UnmarshalPayload] legacy mode failed: %v", err)
		return false
	}
	return true
}
