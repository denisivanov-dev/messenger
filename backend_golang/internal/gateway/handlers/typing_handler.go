package handlers

import (
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
	"messenger/backend_golang/internal/modules/events/typing"
)

func RegisterTyping() {
	hub.Register("event_typing", handleTyping)
}

func handleTyping(c types.ClientLike, env common.Envelope) {
	typing.HandleTypingEvent(c, env)

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(
		rdb,
		c.ID(),
		env.ChatType,
		env.TargetID,
	)
	if !ok {
		c.SendError("access denied")
		return
	}

	c.JoinRoomIfNotJoined(roomID)

	out := typing.BuildTypingEnvelope(
		env.ChatType,
		roomID,
		env.TargetID,
		c.ID(),
		c.Username(),
	)
	c.BroadcastJSON(roomID, out)
}
