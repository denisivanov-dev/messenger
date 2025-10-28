package ws

import (
	"log"
	"net/http"

	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/gateway/client"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	"messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/online"
)

func ServeWS(h *hub.Hub, rdb *rds.Client, w http.ResponseWriter, r *http.Request) {
	conn, userID, username, err := utils.UpgradeAndAuth(w, r)
	if err != nil {
		return
	}

	// Fetch all chat IDs the user is part of
	chatIDs, err := rdb.SMembers(r.Context(), "user:"+userID+":chats").Result()
	if err != nil {
		log.Printf("redis SMembers error: %v", err)
		return
	}

	c := &client.Client{
		Hub:      h,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
		UserID:   userID,
		Name:     username,
		RDB:      rdb,
		Rooms:    make(map[string]struct{}, len(chatIDs)+1),
	}

	// Register user in all their chat rooms
	for _, id := range chatIDs {
		c.Rooms[id] = struct{}{}
	}

	// Also subscribe to the system room for global events like online/offline
	c.Rooms[types.SystemRoom] = struct{}{}

	// Update presence in Redis ===
	if err := online.SetStatus(r.Context(), rdb, userID, "online"); err != nil {
		log.Printf("[ws] SetStatus error: %v", err)
	}

	// Broadcast presence to all clients ===
	h.BroadcastMessage(types.RoomMessage{
		RoomID: types.SystemRoom,
		Data:   online.BuildStatusMessage(userID, "online"),
	})

	h.RegisterClient(c)
	go c.WritePump()
	go c.ReadPump()
}