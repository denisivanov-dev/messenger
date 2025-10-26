package client

import "messenger/backend_golang/internal/gateway/types"

func (c *Client) JoinRoomIfNotJoined(roomID string) {
	if _, ok := c.Rooms[roomID]; ok {
		return
	}
	c.Rooms[roomID] = struct{}{}

	c.Hub.JoinRoom(types.JoinRequest{
		Client: c,
		RoomID: roomID,
	})
}

func (c *Client) LeaveAllExcept(allowedRoomIDs ...string) {
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