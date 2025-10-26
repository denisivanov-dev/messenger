package handlers

import (
	"encoding/json"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
)

func RegisterTyping() {
	hub.Register("typing", handleTyping)
}

func handleTyping(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingInitPrivate
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid typing payload") {
		return
	}

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(rdb, c.ID(), payload.ChatType, payload.ReceiverID)
	if !ok {
		c.SendError("access denied")
		return
	}

	c.JoinRoomIfNotJoined(roomID)
	c.BroadcastJSON(roomID, common.TypingMessage{
		Type:     "typing",
		UserID:   c.ID(),
		Username: c.Username(),
		ChatID:   roomID,
	})
}