package handlers

import (
	"log"
	"github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
)

func RegisterInit() {
	hub.Register("event_init_chat", handleInitChat)
}

// --- HANDLER ---

func handleInitChat(c types.ClientLike, env common.Envelope) {
	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	switch env.ChatType {
	case "global":
		handleGlobalInit(c, env, rdb)
	case "private":
		handlePrivateInit(c, env, rdb)
	default:
		c.SendError("unknown chat_type")
	}
}

// --- INTERNAL FUNCTIONS ---

func handleGlobalInit(c types.ClientLike, env common.Envelope, rdb *redis.Client) {
	roomID := env.ChatID
	if roomID == "" {
		roomID = "1"
	}

	c.LeaveAllExcept(types.SystemRoom, roomID)
	c.JoinRoomIfNotJoined(roomID)

	sendChan := gwutils.GetSendChan(c)
	chat.SendHistory(rdb, roomID, sendChan, 50)

	log.Printf("[init_chat:global] user %s joined room %s", c.ID(), roomID)
}

func handlePrivateInit(c types.ClientLike, env common.Envelope, rdb *redis.Client) {
	roomID, ok := chat.ResolveRoom(rdb, c.ID(), "private", env.TargetID)
	if !ok {
		c.SendError("access denied or invalid chat")
		return
	}

	c.LeaveAllExcept(types.SystemRoom, roomID)
	c.JoinRoomIfNotJoined(roomID)

	sendChan := gwutils.GetSendChan(c)
	chat.SendHistory(rdb, roomID, sendChan, 50)

	log.Printf("[init_chat:private] user %s joined private room %s (target=%s)", c.ID(), roomID, env.TargetID)
}