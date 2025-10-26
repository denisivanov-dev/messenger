package hub

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/gateway/types"
)

// SendToUser sends a payload to a single connected user
func (h *Hub) SendToUser(userID string, payload any) {
	client, ok := h.userClients[userID]
	log.Printf("[hub] SendToUser → %s (connected=%v)", userID, ok)

	if !ok {
		log.Printf("[hub] user %s not connected", userID)
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[hub] marshal error for user %s: %v", userID, err)
		return
	}

	SafeSend(client, data, func() {
		h.Unregister <- client
	})
}

// SendToUsers sends a payload to multiple connected users
func (h *Hub) SendToUsers(userIDs []string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[hub] marshal error in SendToUsers: %v", err)
		return
	}

	for _, uid := range userIDs {
		if client, ok := h.userClients[uid]; ok {
			SafeSend(client, data, func() {
				h.Unregister <- client
			})
		}
	}
}

// SafeSend tries to send data to client and unregisters on failure
func SafeSend(c types.ClientLike, data []byte, onDrop func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in SafeSend: %v", r)
		}
	}()
	c.Send(data)
}