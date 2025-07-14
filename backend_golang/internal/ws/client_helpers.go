package ws

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/chat"
	"messenger/backend_golang/internal/utils"
)

func (c *Client) joinRoomIfNotJoined(roomID string) {
	if _, ok := c.Rooms[roomID]; !ok {
		c.Rooms[roomID] = struct{}{}
		c.Hub.JoinRoom <- joinReq{Client: c, RoomID: roomID}
	}
}

func (c *Client) sendHistory(roomID string, limit int) {
	msgs, err := chat.LoadMessageHistory(c.RDB, roomID, int64(limit))
	if err != nil {
		log.Printf("failed to load chat history for room %s: %v", roomID, err)
		return
	}

	for _, msg := range msgs {
		out, err := json.Marshal(msg)
		if err != nil {
			log.Printf("marshal error: %v", err)
			continue
		}
		c.Send <- out
	}
}

func (c *Client) getRoomID(chatType, receiverID string) string {
	if chatType == "private" {
		return utils.GeneratePrivateChatKey(c.UserID, receiverID)
	}
	return "1"
}

func (c *Client) broadcastJSON(roomID string, payload any) {
	out, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	c.Hub.Broadcast <- RoomMessage{
		RoomID: roomID,
		Data:   out,
	}
}

func (c *Client) leaveAllExcept(allowedRoomIDs ...string) {
	keep := make(map[string]struct{}, len(allowedRoomIDs))
	for _, id := range allowedRoomIDs {
		keep[id] = struct{}{}
	}

	for room := range c.Rooms {
		if _, shouldKeep := keep[room]; !shouldKeep {
			delete(c.Rooms, room)
		}
	}
}
