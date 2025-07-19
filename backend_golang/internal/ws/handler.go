package ws

import (
	"log"
	"net/http"
	"sync"

	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/online"
	"messenger/backend_golang/internal/utils"
)

func ServeWS(hub *Hub, rdb *rds.Client, w http.ResponseWriter, r *http.Request) {
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

	client := &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
		RDB:      rdb,
		Rooms:    make(map[string]struct{}, len(chatIDs)+1),
		once:     sync.Once{},
	}

	// Register user in all their chat rooms
	for _, id := range chatIDs {
		client.Rooms[id] = struct{}{}
	}

	// Also subscribe to the system room for global events like online/offline
	client.Rooms[systemRoom] = struct{}{}

	if err := online.SetOnline(r.Context(), rdb, userID); err != nil {
		log.Printf("online.SetOnline: %v", err)
	}

	hub.Register <- client

	hub.Broadcast <- RoomMessage{
		RoomID: systemRoom,
		Data:   online.BuildStatusMessage(userID, common.Online),
	}

	go client.WritePump()
	go client.ReadPump()
}