package client

import (
	"encoding/json"
	"log"
)

func (c *Client) BroadcastJSON(roomID string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	c.Hub.Broadcast <- RoomMessage{
		RoomID: roomID,
		Data:   data,
	}
}

func (c *Client) SendError(message string) {
	log.Printf("[client %s] %s", c.UserID, message)
	data, _ := json.Marshal(map[string]string{
		"type":    "error",
		"message": message,
	})
	c.Send <- data
}