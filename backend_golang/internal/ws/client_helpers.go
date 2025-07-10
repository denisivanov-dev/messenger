package ws

import (
	"log"

	"messenger/backend_golang/internal/chat"
)

func (c *Client) joinRoomIfNotJoined(roomID string) {
	if _, ok := c.Rooms[roomID]; !ok {
		c.Rooms[roomID] = struct{}{}
		c.Hub.JoinRoom <- joinReq{Client: c, RoomID: roomID}
	}
}

func (c *Client) sendHistory(roomID string, limit int) {
	hist, err := chat.LoadHistoryJSON(c.RDB, roomID, limit)
	if err != nil {
		log.Printf("failed to load chat history for room %s: %v", roomID, err)
		return
	}

	for _, msg := range hist {
		c.Send <- msg
	}
}
