package handlers

import (
	"log"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
)

func RegisterInit() {
	hub.Register("event_init_global", handleInitGlobal)
	hub.Register("event_init_private", handleInitPrivate)
}

// --- HANDLERS ---

func handleInitGlobal(c types.ClientLike, env common.Envelope) {
	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}
	
	roomID := env.ChatID
	if roomID == "" {
		roomID = "1"
	}

	c.LeaveAllExcept(types.SystemRoom, roomID)
	c.JoinRoomIfNotJoined(roomID)

	sendChan := gwutils.GetSendChan(c)
	chat.SendHistory(rdb, roomID, sendChan, 50)

	log.Printf("[init_global] user %s joined room %s", c.ID(), roomID)
}

func handleInitPrivate(c types.ClientLike, env common.Envelope) {
	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(rdb, c.ID(), env.ChatType, env.TargetID)
	if !ok {
		c.SendError("access denied or invalid chat")
		return
	}

	c.LeaveAllExcept(types.SystemRoom, roomID)
	c.JoinRoomIfNotJoined(roomID)

	sendChan := gwutils.GetSendChan(c)
	chat.SendHistory(rdb, roomID, sendChan, 50)

	log.Printf("[init_private] user %s joined private room %s (target=%s)", c.ID(), roomID, env.TargetID)
}