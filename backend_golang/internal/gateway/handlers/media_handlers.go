package handlers

import (
	"encoding/json"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
	"messenger/backend_golang/internal/modules/voice"
)

func RegisterMedia() {
	hub.Register("camera_status", handleCameraStatus)
	hub.Register("screen_status", handleScreenStatus)
	hub.Register("mic_status", handleMicStatus)
}

func handleCameraStatus(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingCameraStatus
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid camera_status payload") {
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

	voice.SetCallMediaStatus(rdb, roomID, c.ID(), "cam", payload.Enabled)
	c.GetHub().SendToUser(payload.ReceiverID, common.OutgoingCameraStatus{
		Type:     "incoming_camera_status",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
		Enabled:  payload.Enabled,
	})
}

func handleScreenStatus(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingScreenStatus
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid screen_status payload") {
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

	voice.SetCallMediaStatus(rdb, roomID, c.ID(), "screen", payload.Enabled)
	c.GetHub().SendToUser(payload.ReceiverID, common.OutgoingScreenStatus{
		Type:     "incoming_screen_status",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
		Enabled:  payload.Enabled,
	})
}

func handleMicStatus(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingMicStatus
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid mic_status payload") {
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

	voice.SetCallMediaStatus(rdb, roomID, c.ID(), "mic", payload.Enabled)
	c.GetHub().SendToUser(payload.ReceiverID, common.OutgoingMicStatus{
		Type:     "incoming_mic_status",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
		Enabled:  payload.Enabled,
	})
}