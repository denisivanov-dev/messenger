package ws

import (
	"context"
	"log"
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"
	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/chat"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/online"
	"messenger/backend_golang/internal/utils"
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
}

func (c *Client) ReadPump() {
	log.Printf("client rooms: %v", c.Rooms)

	defer func() {
		// Notify others that user went offline
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

		switch initMsg.Type {
		case "init_global":
			roomID := "1"
			c.joinRoomIfNotJoined(roomID)
			c.sendHistory(roomID, 50)
			continue

		case "init_private":
			roomID := utils.GeneratePrivateChatKey(c.UserID, initMsg.ReceiverID)
			c.joinRoomIfNotJoined(roomID)
			log.Printf("joined room %s", roomID)
			c.sendHistory(roomID, 50)
			continue

		case "typing":
			roomID := "1"
			if initMsg.ChatType == "private" {
				roomID = utils.GeneratePrivateChatKey(c.UserID, initMsg.ReceiverID)
			}
			c.joinRoomIfNotJoined(roomID)

			typingPayload := map[string]any{
				"type":     "typing",
				"user_id":  c.UserID,
				"username": c.Username,
				"chat_id":  roomID,
			}

			out, err := json.Marshal(typingPayload)
			if err != nil {
				log.Printf("typing marshal error: %v", err)
				break
			}

			c.Hub.Broadcast <- RoomMessage{
				RoomID: roomID,
				Data:   out,
			}
			continue
		}

		chatID, outBytes, ok := chat.HandleIncoming(raw, c.UserID, c.Username, c.RDB)
		if !ok {
			continue
		}

		c.joinRoomIfNotJoined(chatID)

		c.Hub.Broadcast <- RoomMessage{
			RoomID: chatID,
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
			if !ok || c.Conn.WriteMessage(websocket.TextMessage, msg) != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if c.Conn.WriteMessage(websocket.PingMessage, nil) != nil {
				return
			}
		}
	}
}