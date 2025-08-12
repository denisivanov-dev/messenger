package ws

import (
	"encoding/json"
	"log"
)

const systemRoom = "sys"

type RoomMessage struct {
	RoomID string
	Data   []byte
}

type joinReq struct {
	Client *Client
	RoomID string
}

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan RoomMessage
	JoinRoom   chan joinReq

	clients     map[*Client]bool
	userClients map[string]*Client
	rooms       map[string]map[*Client]bool
}

func NewHub() *Hub {
	h := &Hub{
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Broadcast:    make(chan RoomMessage, 256),
		JoinRoom:     make(chan joinReq, 64),
		clients:      make(map[*Client]bool),
		userClients:  make(map[string]*Client),
		rooms:        make(map[string]map[*Client]bool),
	}
	go h.Run()
	return h
}

func (h *Hub) Run() {
	safeSend := func(c *Client, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in safeSend: %v", r)
			}
		}()
		select {
		case c.Send <- data:
		default:
			// Disconnect unresponsive client
			c.once.Do(func() {
				h.Unregister <- c
			})
		}
	}

	for {
		select {
		case c := <-h.Register:
			h.clients[c] = true
			h.userClients[c.UserID] = c 
			for rid := range c.Rooms {
				if h.rooms[rid] == nil {
					h.rooms[rid] = make(map[*Client]bool)
				}
				h.rooms[rid][c] = true
			}
			// Лог всех юзеров после каждого подключения
			log.Printf("Active clients:")
			for uid := range h.userClients {
				log.Printf("• %s", uid)
			}

		case c := <-h.Unregister:
			delete(h.clients, c)
			delete(h.userClients, c.UserID)
			for rid := range c.Rooms {
				if set, ok := h.rooms[rid]; ok && set != nil {
					delete(set, c)
					if len(set) == 0 {
						delete(h.rooms, rid)
					}
				}
			}
			close(c.Send)

		// Add client to a new room dynamically
		case jr := <-h.JoinRoom:
			if h.rooms[jr.RoomID] == nil {
				h.rooms[jr.RoomID] = make(map[*Client]bool)
			}
			h.rooms[jr.RoomID][jr.Client] = true
			jr.Client.Rooms[jr.RoomID] = struct{}{}

		// Route a message to a room
		case msg := <-h.Broadcast:
			if msg.RoomID == systemRoom {
				for c := range h.clients {
					safeSend(c, msg.Data)
				}
				continue
			}
			if set, ok := h.rooms[msg.RoomID]; ok {
				for c := range set {
					if _, alive := h.clients[c]; !alive {
						delete(set, c)
						continue
					}
					safeSend(c, msg.Data)
				}
			}
		}
	}
}

// SendToUser sends payload to a specific user if connected
func (h *Hub) SendToUser(userID string, payload any) {
	client, ok := h.userClients[userID]
	log.Printf("📨 SendToUser: to=%s, connected=%v, payload=%+v", userID, ok, payload)

	if !ok {
		log.Printf("⚠️ user %s not connected", userID)
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Marshal error in SendToUser: %v", err)
		return
	}

	select {
	case client.Send <- data:
		log.Printf("✅ Sent WS payload to user=%s", userID)
	default:
		log.Printf("🚫 Send channel full, unregistering user=%s", userID)
		client.once.Do(func() {
			h.Unregister <- client
		})
	}
}
