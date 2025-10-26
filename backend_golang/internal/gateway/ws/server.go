package ws

import (
	"log"
	"net/http"

	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/client"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	"messenger/backend_golang/internal/modules/online"
	"messenger/backend_golang/internal/gateway/utils"
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

	if err := online.SetOnline(r.Context(), rdb, userID); err != nil {
		log.Printf("online.SetOnline: %v", err)
	}

	h.RegisterClient(c)

	h.BroadcastMessage(types.RoomMessage{
		RoomID: types.SystemRoom,
		Data:   online.BuildStatusMessage(userID, common.Online),
	})

	go c.WritePump()
	go c.ReadPump()
}