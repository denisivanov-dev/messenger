package client

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/gateway/types"
)

func (c *Client) BroadcastJSON(roomID string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[client %s] marshal error: %v", c.UserID, err)
		return
	}

	c.Hub.BroadcastMessage(types.RoomMessage{
		RoomID: roomID,
		Data:   data,
	})
}

func (c *Client) SendError(message string) {
	log.Printf("[client %s] %s", c.UserID, message)

	data, _ := json.Marshal(map[string]string{
		"type":    "error",
		"message": message,
	})

	c.Send(data)
}