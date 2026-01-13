package handlers

import (
	"log"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
	jutils "messenger/backend_golang/internal/utils"
)

func RegisterChat() {
	hub.Register("message_send", handleSendMessage)
	hub.Register("message_delete", handleDeleteMessage)
	hub.Register("message_edit", handleEditMessage)
	hub.Register("message_pin", handlePinMessage)
}

func handleSendMessage(c types.ClientLike, env common.Envelope) {
	var payload common.MessagePayload
	if err := jutils.MapToStruct(env.Payload, &payload); err != nil {
		c.SendError("invalid message payload")
		return
	}

	if payload.Text == "" {
		c.SendError("empty message")
		return
	}
	env.Payload = payload

	if len(payload.Attachments) > 5 {
		c.SendError("too many attachments: max 5")
		return
	}

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(rdb, c.ID(), env.ChatType, env.TargetID)
	if !ok {
		c.SendError("access denied")
		return
	}

	senderID := c.ID()
	username := c.Username()
	outMsg, ok := chat.HandleSendMessage(env, rdb, senderID, username)
	if !ok {
		return
	}

	c.JoinRoomIfNotJoined(roomID)
	c.BroadcastJSON(roomID, outMsg)

	log.Printf("[message_send] %s sent message in room %s", c.ID(), roomID)
}

func handleDeleteMessage(c types.ClientLike, env common.Envelope) {
	var payload common.MessagePayload
	jutils.MapToStruct(env.Payload, &payload)
	env.Payload = payload

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(rdb, c.ID(), env.ChatType, env.TargetID)
	if !ok {
		c.SendError("access denied")
		return
	}

	c.JoinRoomIfNotJoined(roomID)

	if deleted := chat.DeleteMessageFromRedisHistory(rdb, roomID, payload.MessageID, c.ID()); deleted != nil {
		c.BroadcastJSON(roomID, deleted)
		log.Printf("[message_delete] %s deleted message in room %s", c.ID(), roomID)
	}
}

func handleEditMessage(c types.ClientLike, env common.Envelope) {
	var payload common.MessagePayload
	jutils.MapToStruct(env.Payload, &payload)
	env.Payload = payload

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(rdb, c.ID(), env.ChatType, env.TargetID)
	if !ok {
		c.SendError("access denied")
		return
	}

	c.JoinRoomIfNotJoined(roomID)

	if edited := chat.EditMessageInRedisHistory(rdb, roomID, payload.MessageID, payload.Text, c.ID()); edited != nil {
		c.BroadcastJSON(roomID, edited)
		log.Printf("[message_edit] %s edited message in room %s", c.ID(), roomID)
	}
}

func handlePinMessage(c types.ClientLike, env common.Envelope) {
	var payload common.MessagePayload
	jutils.MapToStruct(env.Payload, &payload)
	env.Payload = payload

	rdb := gwutils.GetRedis(c)
	if rdb == nil {
		c.SendError("redis unavailable")
		return
	}

	roomID, ok := chat.ResolveRoom(rdb, c.ID(), env.ChatType, env.TargetID)
	if !ok {
		c.SendError("access denied")
		return
	}

	c.JoinRoomIfNotJoined(roomID)

	pin := payload.Pinned
	if pinned := chat.PinMessageInRedisHistory(rdb, roomID, payload.MessageID, pin, c.ID()); pinned != nil {
		c.BroadcastJSON(roomID, pinned)
		log.Printf("[message_pin] %s toggled pin=%t in room %s", c.ID(), pin, roomID)
	}
}