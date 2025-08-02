package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/chat"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/online"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Username string
	RDB      *rds.Client
	Rooms    map[string]struct{}
	once     sync.Once
}

func (c *Client) ReadPump() {
	log.Printf("client rooms: %v", c.Rooms)

	defer func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recover from panic in defer: %v", r)
			}
		}()
		c.Hub.Broadcast <- RoomMessage{
			RoomID: systemRoom,
			Data:   online.BuildStatusMessage(c.UserID, common.Offline),
		}
		_ = online.SetOffline(context.Background(), c.RDB, c.UserID)
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	c.prepareConn()

	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			return
		}

		var msgType common.IncomingMessage
		if err := json.Unmarshal(raw, &msgType); err != nil {
			log.Printf("unmarshal type error: %v", err)
			continue
		}

		switch msgType.Type {

		case "init_global":
			var payload common.IncomingInitGlobal
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid init_global payload")
				continue
			}
			roomID := "1"
			c.leaveAllExcept("sys", roomID)
			c.joinRoomIfNotJoined(roomID)
			c.sendHistory(roomID, 50)

		case "init_private":
			var payload common.IncomingInitPrivate
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid init_private payload")
				continue
			}
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied or invalid chat")
				continue
			}
			c.leaveAllExcept("sys", roomID)
			c.joinRoomIfNotJoined(roomID)
			c.sendHistory(roomID, 50)

		case "typing":
			var payload common.IncomingInitPrivate
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid typing payload")
				continue
			}
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			c.broadcastJSON(roomID, common.TypingMessage{
				Type:     "typing",
				UserID:   c.UserID,
				Username: c.Username,
				ChatID:   roomID,
			})

		case "send_message":
			var payload common.IncomingSendMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid send_message payload")
				continue
			}

			if len(payload.Attachments) > 5 {
				c.sendError("too many attachments: max 5")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue 
			}

			// log.Printf("📨 Incoming send_message: user=%s attachments=%d", c.UserID, len(payload.Attachments))
			
			outMsg, ok := chat.HandleSendMessage(payload, c.UserID, c.Username, c.RDB)
			if !ok {
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			c.broadcastJSON(roomID, outMsg)

		case "delete_message":
			var payload common.IncomingDeleteMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid delete_message payload")
				continue
			}
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			if deleted := chat.DeleteMessageFromRedisHistory(c.RDB, roomID, payload.MessageID, c.UserID); deleted != nil {
				c.broadcastJSON(roomID, deleted)
			}

		case "edit_message":
			var payload common.IncomingEditMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid edit_message payload")
				continue
			}
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			if edited := chat.EditMessageInRedisHistory(c.RDB, roomID, payload.MessageID, payload.NewText, c.UserID); edited != nil {
				c.broadcastJSON(roomID, edited)
			}

		case "pin_message":
			var payload common.IncomingPinMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid edit_message payload")
				continue
			}

			log.Printf(payload.ChatType)
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			pin := payload.Action == "pin"

			c.joinRoomIfNotJoined(roomID)
			if pinned := chat.PinMessageInRedisHistory(c.RDB, roomID, payload.MessageID, pin); pinned != nil {
				c.broadcastJSON(roomID, pinned)
			}

		default:
			c.sendError("unsupported message type: " + msgType.Type)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			log.Printf("client rooms: %v", c.Rooms)
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}