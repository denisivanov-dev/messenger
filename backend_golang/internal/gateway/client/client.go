package client

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	"messenger/backend_golang/internal/modules/online"
)

const (
	WriteWait  = 10 * time.Second
	PongWait   = 60 * time.Second
	PingPeriod = PongWait * 9 / 10
)

type Client struct {
	Hub      types.HubLike
	Conn     *websocket.Conn
	SendChan chan []byte
	UserID   string
	Name     string
	RDB      *rds.Client
	Rooms    map[string]struct{}

	once     sync.Once
}

// === ClientLike Interface Implementation ===

func (c *Client) ID() string { 
	return c.UserID 
}

func (c *Client) Username() string {
	return c.Name
}

func (c *Client) GetHub() types.HubLike {
	return c.Hub
}

func (c *Client) GetRedis() *rds.Client {
	return c.RDB
}

func (c *Client) GetSendChannel() chan []byte {
	return c.SendChan
}

func (c *Client) GetRooms() map[string]struct{} {
	return c.Rooms
}

func (c *Client) Send(data []byte) {
	select {
	case c.SendChan <- data:
	default:
		log.Printf("[client %s] send buffer full", c.UserID)
	}
}

// --- Core network pumps ---

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
		// _ = raw // placeholder to avoid unused var
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
		case msg, ok := <-c.SendChan:
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

// --- Cleanup on disconnect ---

func (c *Client) cleanup() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[client %s] recovered in cleanup: %v", c.UserID, r)
		}
	}()

	log.Printf("[client %s] cleanup started", c.UserID)

	// Notify all clients about status change (broadcast envelope)
	statusMsg := online.BuildStatusMessage(c.UserID, "offline")
	c.BroadcastJSON(types.SystemRoom, statusMsg)

	// Update Redis presence
	ctx := context.Background()
	_ = online.SetStatus(ctx, c.RDB, c.UserID, "offline")

	c.Hub.UnregisterClient(c)
	_ = c.Conn.Close()

	log.Printf("[client %s] disconnected", c.UserID)
}