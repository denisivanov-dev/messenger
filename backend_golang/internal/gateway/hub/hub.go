package hub

import (
	"messenger/backend_golang/internal/gateway/types"
)

type Hub struct {
	Register   chan types.ClientLike
	Unregister chan types.ClientLike
	Broadcast  chan types.RoomMessage
	JoinRoomChan chan types.JoinRequest


	clients     map[types.ClientLike]bool
	userClients map[string]types.ClientLike
	rooms       map[string]map[types.ClientLike]bool
}

func NewHub() *Hub {
	h := &Hub{
		Register:     make(chan types.ClientLike),
		Unregister:   make(chan types.ClientLike),
		Broadcast:    make(chan types.RoomMessage, 256),
		JoinRoomChan:     make(chan types.JoinRequest, 64),
		clients:      make(map[types.ClientLike]bool),
		userClients:  make(map[string]types.ClientLike),
		rooms:        make(map[string]map[types.ClientLike]bool),
	}
	go h.Run()
	return h
}

// === HubLike Interface ===

func (h *Hub) BroadcastMessage(msg types.RoomMessage) {
	h.Broadcast <- msg
}

func (h *Hub) JoinRoom(req types.JoinRequest) {
	h.JoinRoomChan <- req
}

func (h *Hub) UnregisterClient(c types.ClientLike) {
	h.Unregister <- c
}

func (h *Hub) RegisterClient(c types.ClientLike) {
	h.Register <- c
}


func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.handleRegister(c)

		case c := <-h.Unregister:
			h.handleUnregister(c)

		// Add client to a new room dynamically
		case jr := <-h.JoinRoomChan:
			h.handleJoinRoom(jr)

		// Route a message to a room
		case msg := <-h.Broadcast:
			h.handleBroadcast(msg)
		}
	}
}