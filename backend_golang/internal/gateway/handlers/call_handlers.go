package handlers

import (
	"encoding/json"
	"log"
	"time"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
	"messenger/backend_golang/internal/modules/chat"
	"messenger/backend_golang/internal/modules/voice"
)

func RegisterCalls() {
	hub.Register("start_call", handleStartCall)
	hub.Register("cancel_call", handleCancelCall)
	hub.Register("call_answer", handleCallAnswer)
	hub.Register("join_call", handleJoinCall)
	hub.Register("leave_call", handleLeaveCall)
}

// --- helpers ---

// sendToUser safely sends a payload to a specific user via hub.
func sendToUser(c types.ClientLike, userID string, payload any) {
	if h := gwutils.GetHub(c); h != nil {
		h.SendToUser(userID, payload)
	}
}

// --- HANDLERS ---

func handleStartCall(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingStartCall
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid start_call payload") {
		return
	}

	if payload.ChatType != "private" || payload.ReceiverID == "" {
		c.SendError("invalid call context")
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

	voice.SetCallParticipant(rdb, roomID, c.ID(), "joined")
	voice.SetCallParticipant(rdb, roomID, payload.ReceiverID, "calling")
	voice.AddCallHistoryParticipant(rdb, roomID, c.ID())

	sendToUser(c, payload.ReceiverID, common.OutgoingCallNotification{
		Type:     "incoming_call",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
	})

	callInfo := &common.CallInfo{
		Status:       "ongoing",
		Participants: []string{c.ID()},
		StartedAt:    time.Now().Unix(),
	}

	msg, ok := chat.HandleSystemMessage(
		"call_started",
		payload.ChatType,
		c.ID(),
		payload.ReceiverID,
		callInfo,
		rdb,
	)
	if ok {
		c.JoinRoomIfNotJoined(msg.ChatID)
		c.BroadcastJSON(msg.ChatID, msg)
	}
}

func handleCancelCall(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingCancelCall
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid cancel_call payload") {
		return
	}

	if payload.ChatType != "private" || payload.ReceiverID == "" {
		c.SendError("invalid call context")
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

	voice.RemoveCallParticipant(rdb, roomID, payload.ReceiverID)

	sendToUser(c, payload.ReceiverID, common.OutgoingCancelCallNotification{
		Type:     "incoming_cancel_call",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
	})
}

func handleCallAnswer(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingCallAnswer
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid call_answer payload") {
		return
	}

	if payload.ChatType != "private" || payload.ReceiverID == "" {
		c.SendError("invalid call context")
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

	if payload.Accepted {
		voice.SetCallParticipant(rdb, roomID, c.ID(), "joined")
	} else {
		voice.RemoveCallParticipant(rdb, roomID, c.ID())
	}

	sendToUser(c, payload.ReceiverID, common.OutgoingCallAnswer{
		Type:     "incoming_call_answer",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
		Accepted: payload.Accepted,
	})
}

func handleJoinCall(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingJoinCall
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid join_call payload") {
		return
	}

	if payload.ChatType != "private" || payload.ReceiverID == "" {
		c.SendError("invalid call context")
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

	voice.SetCallParticipant(rdb, roomID, c.ID(), "joined")
	voice.AddCallHistoryParticipant(rdb, roomID, c.ID())

	sendToUser(c, payload.ReceiverID, common.OutgoingJoinCallNotification{
		Type:     "incoming_join_call",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
	})
}

func handleLeaveCall(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingLeaveCall
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid leave_call payload") {
		return
	}

	if payload.ChatType != "private" || payload.ReceiverID == "" {
		c.SendError("invalid call context")
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

	voice.RemoveCallParticipant(rdb, roomID, c.ID())

	sendToUser(c, payload.ReceiverID, common.OutgoingLeaveCallNotification{
		Type:     "incoming_leave_call",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
	})

	remaining, _ := voice.GetCallParticipants(rdb, roomID)
	if len(remaining) == 0 {
		voice.ClearCallRoom(rdb, roomID)
		msgID, err := voice.GetOngoingCallMessageID(rdb, roomID)
		if err == nil && msgID != "" {
			startedAt, _ := voice.GetCallStartTime(rdb, roomID)
			now := time.Now().Unix()
			duration := now - startedAt

			participants := voice.GetCallHistoryParticipants(rdb, roomID)
			updatedMsg := common.OutgoingMessage{
				Type:       "call_started",
				MessageID:  msgID,
				ChatID:     roomID,
				Timestamp:  startedAt * 1000,
				Username:   "system",
				UserID:     "0",
				ReceiverID: payload.ReceiverID,
				Pinned:     false,
				CallInfo: &common.CallInfo{
					Status:       "ended",
					StartedAt:    startedAt,
					Duration:     duration,
					Participants: participants,
				},
			}

			chat.UpdateSystemMessageInRedisHistory(rdb, roomID, msgID, updatedMsg)
			c.BroadcastJSON(roomID, updatedMsg)

			voice.ClearOngoingCallMessageID(rdb, roomID)
			voice.ClearCallStartTime(rdb, roomID)
			voice.ClearCallHistoryParticipants(rdb, roomID)
		}
	}

	log.Printf("[leave_call] %s left call in %s", c.ID(), roomID)
}