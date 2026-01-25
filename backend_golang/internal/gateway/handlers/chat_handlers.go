package handlers

import (
	"log"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
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
	if len(payload.Attachments) > 5 {
		c.SendError("too many attachments: max 5")
		return
	}

	env.Payload = payload

	ctx, err := chat.ResolveRoomContext(c, env)
	if err != nil {
		c.SendError(err.Error())
		return
	}

	outMsg, ok := chat.HandleSendUserMessage(
		env,
		ctx.RoomID,
		ctx.RDB,
		c.ID(),
		c.Username(),
	)
	if !ok {
		return
	}

	c.BroadcastJSON(ctx.RoomID, outMsg)
	log.Printf("[message_send] %s sent message in room %s", c.ID(), ctx.RoomID)
}

func handleDeleteMessage(c types.ClientLike, env common.Envelope) {
	var payload common.MessagePayload
	if err := jutils.MapToStruct(env.Payload, &payload); err != nil {
		c.SendError("invalid message payload")
		return
	}

	env.Payload = payload

	ctx, err := chat.ResolveRoomContext(c, env)
	if err != nil {
		c.SendError(err.Error())
		return
	}

	deleted := chat.DeleteMessageFromRedisHistory(ctx.RDB, ctx.RoomID, payload.MessageID, c.ID(), c.Username())
	if deleted == nil {
		return
	}

	c.BroadcastJSON(ctx.RoomID, deleted)
	log.Printf("[message_delete] %s deleted message in room %s", c.ID(), ctx.RoomID)
}

func handleEditMessage(c types.ClientLike, env common.Envelope) {
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

	ctx, err := chat.ResolveRoomContext(c, env)
	if err != nil {
		c.SendError(err.Error())
		return
	}

	edited := chat.EditMessageInRedisHistory(ctx.RDB, ctx.RoomID, payload.MessageID, payload.Text, c.ID(), c.Username())
	if edited == nil {
		return
	}

	c.BroadcastJSON(ctx.RoomID, edited)
	log.Printf("[message_edit] %s edited message in room %s", c.ID(), ctx.RoomID)
}

func handlePinMessage(c types.ClientLike, env common.Envelope) {
	var payload common.MessagePayload
	if err := jutils.MapToStruct(env.Payload, &payload); err != nil {
		c.SendError("invalid message payload")
		return
	}
	
	env.Payload = payload

	ctx, err := chat.ResolveRoomContext(c, env)
	if err != nil {
		c.SendError(err.Error())
		return
	}

	pinned := chat.PinMessageInRedisHistory(ctx.RDB, ctx.RoomID, payload.MessageID, payload.Pinned, c.ID())
	if pinned == nil {
		return
	}

	c.BroadcastJSON(ctx.RoomID, pinned)
	log.Printf("[message_pin] %s toggled pin=%t in room %s", c.ID(), payload.Pinned, ctx.RoomID)
}
