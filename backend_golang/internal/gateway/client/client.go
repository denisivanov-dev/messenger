package client

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/online"
)

type Client struct {
	Hub      *hub.Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Username string
	RDB      *rds.Client
	Rooms    map[string]struct{}
	once     sync.Once
}

func (c *Client) ReadPump() {
	log.Printf("[client %s] connected", c.UserID)
	defer c.cleanup()

	c.prepareConn()

	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[client %s] read error: %v", c.UserID, err)
			}
			return
		}

		hub.Handle(c, raw)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				log.Printf("[client %s] send channel closed", c.UserID)
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[client %s] write error: %v", c.UserID, err)
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[client %s] ping failed: %v", c.UserID, err)
				return
			}
		}
	}
}

func (c *Client) cleanup() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[client %s] recovered in cleanup: %v", c.UserID, r)
		}
	}()

	log.Printf("[client %s] cleanup started", c.UserID)

	c.Hub.Broadcast <- hub.RoomMessage{
		RoomID: common.SystemRoomID,
		Data:   online.BuildStatusMessage(c.UserID, common.Offline),
	}

	_ = online.SetOffline(context.Background(), c.RDB, c.UserID)

	c.Hub.Unregister <- c

	_ = c.Conn.Close()

	log.Printf("[client %s] disconnected", c.UserID)
}