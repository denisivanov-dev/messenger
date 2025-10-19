package hub

import "log"
// — Hub's internal functions for managing clients, rooms, and mailing lists.

func (h *Hub) handleRegister(c *Client) {
	h.clients[c] = true
	h.userClients[c.UserID] = c
	for rid := range c.Rooms {
		if h.rooms[rid] == nil {
			h.rooms[rid] = make(map[*Client]bool)
		}
		h.rooms[rid][c] = true
	}

	// info about active users
	log.Printf("Active clients:")
	for uid := range h.userClients {
		log.Printf("• %s", uid)
	}
}

func (h *Hub) handleUnregister(c *Client) {
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
}

func (h *Hub) handleJoinRoom(jr joinReq) {
	if h.rooms[jr.RoomID] == nil {
		h.rooms[jr.RoomID] = make(map[*Client]bool)
	}
	h.rooms[jr.RoomID][jr.Client] = true
	jr.Client.Rooms[jr.RoomID] = struct{}{}
}

func (h *Hub) handleBroadcast(msg RoomMessage) {
	if msg.RoomID == systemRoom {
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