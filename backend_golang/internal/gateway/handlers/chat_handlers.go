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

func RegisterChat() {
	hub.Register("send_message", handleSendMessage)
	hub.Register("delete_message", handleDeleteMessage)
	hub.Register("edit_message", handleEditMessage)
	hub.Register("pin_message", handlePinMessage)
}

// --- HANDLERS ---

func handleSendMessage(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingSendMessage
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid send_message payload") {
		return
	}

	if len(payload.Attachments) > 5 {
		c.SendError("too many attachments: max 5")
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

	outMsg, ok := chat.HandleSendMessage(payload, c.ID(), c.Username(), rdb)
	if !ok {
		return
	}

	c.JoinRoomIfNotJoined(roomID)
	c.BroadcastJSON(roomID, outMsg)
	log.Printf("[send_message] %s sent message in room %s", c.ID(), roomID)
}

func handleDeleteMessage(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingDeleteMessage
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid delete_message payload") {
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
	if deleted := chat.DeleteMessageFromRedisHistory(rdb, roomID, payload.MessageID, c.ID()); deleted != nil {
		c.BroadcastJSON(roomID, deleted)
		log.Printf("[delete_message] %s deleted message in room %s", c.ID(), roomID)
	}
}

func handleEditMessage(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingEditMessage
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid edit_message payload") {
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
	if edited := chat.EditMessageInRedisHistory(rdb, roomID, payload.MessageID, payload.NewText, c.ID()); edited != nil {
		c.BroadcastJSON(roomID, edited)
		log.Printf("[edit_message] %s edited message in room %s", c.ID(), roomID)
	}
}

func handlePinMessage(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingPinMessage
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid pin_message payload") {
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
	pin := payload.Action == "pin"

	if pinned := chat.PinMessageInRedisHistory(rdb, roomID, payload.MessageID, pin, c.ID()); pinned != nil {
		c.BroadcastJSON(roomID, pinned)
		log.Printf("[pin_message] %s toggled pin=%t in room %s", c.ID(), pin, roomID)
	}
}
