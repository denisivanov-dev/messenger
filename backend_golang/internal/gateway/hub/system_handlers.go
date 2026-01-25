package hub

import (
	"log"

	"messenger/backend_golang/internal/gateway/types"
)

// — Hub's internal functions for managing clients, rooms, and mailing lists.

func (h *Hub) handleRegister(c types.ClientLike) {
	h.clients[c] = true
	h.userClients[c.ID()] = c
	for rid := range c.GetRooms() {
		if h.rooms[rid] == nil {
			h.rooms[rid] = make(map[types.ClientLike]bool)
		}
		h.rooms[rid][c] = true
	}

	// info about active users
	log.Printf("Active clients:")
	for uid := range h.userClients {
		log.Printf("• %s", uid)
	}
}

func (h *Hub) handleUnregister(c types.ClientLike) {
	delete(h.clients, c)
	delete(h.userClients, c.ID())
	for rid := range c.GetRooms() {
		if set, ok := h.rooms[rid]; ok && set != nil {
			delete(set, c)
			if len(set) == 0 {
				delete(h.rooms, rid)
			}
		}
	}
}

func (h *Hub) handleJoinRoom(jr types.JoinRequest) {
	if h.rooms[jr.RoomID] == nil {
		h.rooms[jr.RoomID] = make(map[types.ClientLike]bool)
	}
	h.rooms[jr.RoomID][jr.Client] = true
	jr.Client.GetRooms()[jr.RoomID] = struct{}{}
}

func (h *Hub) handleBroadcast(msg types.RoomMessage) {
	if msg.RoomID == types.SystemRoom {
		for c := range h.clients {
			SafeSend(c, msg.Data, func() {
				h.Unregister <- c
			})
		}
		return
	}
	if set, ok := h.rooms[msg.RoomID]; ok {
		for c := range set {
			if _, alive := h.clients[c]; !alive {
				delete(set, c)
				continue
			}
			SafeSend(c, msg.Data, func() {
				h.Unregister <- c
			})
		}
	}
}
