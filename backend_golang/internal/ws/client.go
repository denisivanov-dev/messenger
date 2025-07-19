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

		var initMsg common.IncomingMessage
		if err := json.Unmarshal(raw, &initMsg); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		log.Printf("RECEIVED: %+v", initMsg)

		switch initMsg.Type {
		case "init_global":
			roomID := "1"
			c.leaveAllExcept("sys", roomID)
			c.joinRoomIfNotJoined(roomID)
			c.sendHistory(roomID, 50)
			continue

		case "init_private":
			roomID := c.getRoomID(initMsg.ChatType, initMsg.ReceiverID)
			c.leaveAllExcept("sys", roomID)
			c.joinRoomIfNotJoined(roomID)
			log.Printf("joined room %s", roomID)
			c.sendHistory(roomID, 50)
			continue

		case "typing":
			roomID := c.getRoomID(initMsg.ChatType, initMsg.ReceiverID)
			c.joinRoomIfNotJoined(roomID)
			payload := common.TypingMessage{
				Type:     "typing",
				UserID:   c.UserID,
				Username: c.Username,
				ChatID:   roomID,
			}
			c.broadcastJSON(roomID, payload)
			continue

		case "delete_message":
			roomID := c.getRoomID(initMsg.ChatType, initMsg.ReceiverID)
			c.joinRoomIfNotJoined(roomID)
			if deleted := chat.DeleteMessageFromRedisHistory(c.RDB, roomID, initMsg.MessageID, c.UserID); deleted != nil {
				c.broadcastJSON(roomID, deleted)
			}
			continue

		case "edit_message":
			roomID := c.getRoomID(initMsg.ChatType, initMsg.ReceiverID)
			c.joinRoomIfNotJoined(roomID)
			log.Printf("editing")
			if edited := chat.EditMessageInRedisHistory(c.RDB, roomID, initMsg.MessageID, initMsg.NewText, c.UserID); edited != nil {
				c.broadcastJSON(roomID, edited)
			}
			continue
		}

		roomID, outBytes, ok := chat.HandleIncoming(raw, c.UserID, c.Username, c.RDB)
		if !ok {
			continue
		}
		c.joinRoomIfNotJoined(roomID)

		c.Hub.Broadcast <- RoomMessage{
			RoomID: roomID,
			Data:   outBytes,
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