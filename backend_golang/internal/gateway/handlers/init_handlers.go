package handlers

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
)

func RegisterInit() {
	hub.Register("init_global", handleInitGlobal)
	hub.Register("init_private", handleInitPrivate)
}

// --- HANDLERS ---

func handleInitGlobal(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingInitGlobal
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid init_global payload") {
		return
	}

	roomID := "1" // system global chat room
	c.LeaveAllExcept(types.SystemRoom, roomID)
	c.JoinRoomIfNotJoined(roomID)

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	sendChan := gwutils.GetSendChan(c)
	chat.SendHistory(rdb, roomID, sendChan, 50)

	log.Printf("[init_global] user %s joined room %s", c.ID(), roomID)
}

func handleInitPrivate(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingInitPrivate
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid init_private payload") {
		return
	}

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(rdb, c.ID(), payload.ChatType, payload.ReceiverID)
	if !ok {
		c.SendError("access denied or invalid chat")
		return
	}

	c.LeaveAllExcept(types.SystemRoom, roomID)
	c.JoinRoomIfNotJoined(roomID)

	sendChan := gwutils.GetSendChan(c)
	chat.SendHistory(rdb, roomID, sendChan, 50)

	log.Printf("[init_private] user %s joined private room %s", c.ID(), roomID)
}
