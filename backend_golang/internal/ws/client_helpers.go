package ws

import (
	"encoding/json"
	"log"
	"fmt"
	"context"

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

func (c *Client) resolveRoomID(chatType, receiverID string) (string, bool) {
	ctx := context.Background()

	chatKey := c.getRoomID(chatType, receiverID)
	if chatKey == "1" {
		return "1", true
	}

	chatID, err := c.RDB.Get(ctx, fmt.Sprintf("chat_id:%s", chatKey)).Result()
	if err != nil {
		log.Printf("Redis GET chat_id:%s failed: %v", chatKey, err)
		return "", false
	}

	ok, err := c.RDB.SIsMember(ctx, fmt.Sprintf("chat:%s:participants", chatID), c.UserID).Result()
	if err != nil {
		log.Printf("Redis SISMEMBER failed: %v", err)
		return "", false
	}
	if !ok {
		log.Printf("Access denied: user %s not in chat %s", c.UserID, chatID)
		return "", false
	}

	return chatKey, true
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

func (c *Client) sendError(message string) {
	payload := map[string]string{
		"type":    "error",
		"message": message,
	}
	out, _ := json.Marshal(payload)
	c.Send <- out
}
