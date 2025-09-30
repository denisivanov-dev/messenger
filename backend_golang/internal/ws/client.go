package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/chat"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/online"
	"messenger/backend_golang/internal/voice"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Username string
	RDB      *rds.Client
	Rooms    map[string]struct{}
	once     sync.Once
}

func (c *Client) ReadPump() {
	log.Printf("client rooms: %v", c.Rooms)

	defer func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recover from panic in defer: %v", r)
			}
		}()
		c.Hub.Broadcast <- RoomMessage{
			RoomID: systemRoom,
			Data:   online.BuildStatusMessage(c.UserID, common.Offline),
		}
		_ = online.SetOffline(context.Background(), c.RDB, c.UserID)
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	c.prepareConn()

	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			return
		}

		var msgType common.IncomingMessage
		if err := json.Unmarshal(raw, &msgType); err != nil {
			log.Printf("unmarshal type error: %v", err)
			continue
		}

		switch msgType.Type {

		case "init_global":
			var payload common.IncomingInitGlobal
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid init_global payload")
				continue
			}
			roomID := "1"
			c.leaveAllExcept("sys", roomID)
			c.joinRoomIfNotJoined(roomID)
			c.sendHistory(roomID, 50)

		case "init_private":
			var payload common.IncomingInitPrivate
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid init_private payload")
				continue
			}
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied or invalid chat")
				continue
			}
			c.leaveAllExcept("sys", roomID)
			c.joinRoomIfNotJoined(roomID)
			c.sendHistory(roomID, 50)

		case "typing":
			var payload common.IncomingInitPrivate
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid typing payload")
				continue
			}
			
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			c.broadcastJSON(roomID, common.TypingMessage{
				Type:     "typing",
				UserID:   c.UserID,
				Username: c.Username,
				ChatID:   roomID,
			})

		case "send_message":
			var payload common.IncomingSendMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid send_message payload")
				continue
			}

			if len(payload.Attachments) > 5 {
				c.sendError("too many attachments: max 5")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue 
			}

			// log.Printf("📨 Incoming send_message: user=%s attachments=%d", c.UserID, len(payload.Attachments))
			
			outMsg, ok := chat.HandleSendMessage(payload, c.UserID, c.Username, c.RDB)
			if !ok {
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			c.broadcastJSON(roomID, outMsg)

		case "delete_message":
			var payload common.IncomingDeleteMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid delete_message payload")
				continue
			}
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			if deleted := chat.DeleteMessageFromRedisHistory(c.RDB, roomID, payload.MessageID, c.UserID); deleted != nil {
				c.broadcastJSON(roomID, deleted)
			}

		case "edit_message":
			var payload common.IncomingEditMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid edit_message payload")
				continue
			}
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			c.joinRoomIfNotJoined(roomID)
			if edited := chat.EditMessageInRedisHistory(c.RDB, roomID, payload.MessageID, payload.NewText, c.UserID); edited != nil {
				c.broadcastJSON(roomID, edited)
			}

		case "pin_message":
			var payload common.IncomingPinMessage
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid edit_message payload")
				continue
			}

			log.Printf(payload.ChatType)
			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}
			pin := payload.Action == "pin"

			c.joinRoomIfNotJoined(roomID)

			if pinned := chat.PinMessageInRedisHistory(c.RDB, roomID, payload.MessageID, pin, c.UserID); pinned != nil {
				c.broadcastJSON(roomID, pinned)
			}

		case "start_call":
			var payload common.IncomingStartCall
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid start_call payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			voice.SetCallParticipant(c.RDB, roomID, c.UserID, "joined")
			voice.SetCallParticipant(c.RDB, roomID, payload.ReceiverID, "calling")

			voice.AddCallHistoryParticipant(c.RDB, roomID, c.UserID)

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingCallNotification{
				Type:     "incoming_call",
				FromUser: c.UserID,
				ChatType: payload.ChatType,
			})

			callInfo := &common.CallInfo{
				Status:       "ongoing",
				Participants: []string{c.UserID},
				StartedAt:    time.Now().Unix(),
			}

			msg, ok := chat.HandleSystemMessage(
				"call_started",
				payload.ChatType,
				c.UserID,
				payload.ReceiverID,
				callInfo,
				c.RDB,
			)

			if ok {
				c.joinRoomIfNotJoined(msg.ChatID)
				c.broadcastJSON(msg.ChatID, msg)
			}

		case "cancel_call":
			var payload common.IncomingCancelCall
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid cancel_call payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			voice.RemoveCallParticipant(c.RDB, roomID, payload.ReceiverID)

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingCancelCallNotification{
				Type:     "incoming_cancel_call",
				FromUser: c.UserID,
				ChatType: payload.ChatType,
			})

		case "call_answer":
			var payload common.IncomingCallAnswer
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid call_answer payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			if payload.Accepted {
				voice.SetCallParticipant(c.RDB, roomID, c.UserID, "joined")
			} else {
				voice.RemoveCallParticipant(c.RDB, roomID, c.UserID)
			}

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingCallAnswer{
				Type:     "incoming_call_answer",
				FromUser: c.UserID,
				ChatType: payload.ChatType,
				Accepted: payload.Accepted,
			})

		case "join_call":
			var payload common.IncomingJoinCall
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid join_call payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			voice.SetCallParticipant(c.RDB, roomID, c.UserID, "joined")

			voice.AddCallHistoryParticipant(c.RDB, roomID, c.UserID)

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingJoinCallNotification{
				Type:     "incoming_join_call",
				FromUser: c.UserID,
				ChatType: payload.ChatType,
			})

		case "leave_call":
			var payload common.IncomingLeaveCall
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid leave_call payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			voice.RemoveCallParticipant(c.RDB, roomID, c.UserID)

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingLeaveCallNotification{
				Type:     "incoming_leave_call",
				FromUser: c.UserID,
				ChatType: payload.ChatType,
			})

			remaining, _ := voice.GetCallParticipants(c.RDB, roomID)
			if len(remaining) == 0 {
				voice.ClearCallRoom(c.RDB, roomID)

				msgID, err := voice.GetOngoingCallMessageID(c.RDB, roomID)
				if err == nil && msgID != "" {
					startedAt, _ := voice.GetCallStartTime(c.RDB, roomID)
					now := time.Now().Unix()

					var duration int64 = 0
					if startedAt > 0 {
						duration = now - startedAt
					}

					participants := voice.GetCallHistoryParticipants(c.RDB, roomID)

					updatedMsg := common.OutgoingMessage{
						Type:      "call_started",
						MessageID: msgID,
						ChatID:    roomID,
						CallInfo: &common.CallInfo{
							Status:       "ended",
							StartedAt:    startedAt,
							Duration:     duration,
							Participants: participants,
						},
					}

					chat.UpdateSystemMessageInRedisHistory(c.RDB, roomID, msgID, updatedMsg)
					c.broadcastJSON(roomID, updatedMsg)

					voice.ClearOngoingCallMessageID(c.RDB, roomID)
					voice.ClearCallStartTime(c.RDB, roomID)
					voice.ClearCallHistoryParticipants(c.RDB, roomID)
				}
			}

		case "webrtc_offer":
			var payload common.IncomingWebRTCOffer
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid webrtc_offer payload")
				continue
			}

			if payload.ReceiverID == "" || payload.Offer == nil {
				c.sendError("invalid webrtc_offer structure")
				continue
			}

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingWebRTCOffer{
				Type:     "incoming_webrtc_offer",
				FromUser: c.UserID,
				ChatType: payload.ChatType,
				Offer:    payload.Offer,
			})

		case "webrtc_answer":
			var payload common.IncomingWebRTCAnswer
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid webrtc_answer payload")
				continue
			}

			if payload.ReceiverID == "" || payload.Answer == nil {
				c.sendError("invalid webrtc_answer structure")
				continue
			}

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingWebRTCAnswer{
				Type:     "incoming_webrtc_answer",
				FromUser: c.UserID,
				ChatType: payload.ChatType,
				Answer:   payload.Answer,
			})

		case "ice_candidate":
			var payload common.IncomingIceCandidate
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid ice_candidate payload")
				continue
			}

			if payload.ReceiverID == "" || payload.Candidate == nil {
				c.sendError("invalid ice_candidate structure")
				continue
			}

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingIceCandidate{
				Type:      "incoming_ice_candidate",
				FromUser:  c.UserID,
				ChatType:  payload.ChatType,
				Candidate: payload.Candidate,
			})

		case "camera_status":
			var payload common.IncomingCameraStatus
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid camera_status payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			voice.SetCallMediaStatus(c.RDB, roomID, c.UserID, "cam", payload.Enabled)

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingCameraStatus{
				Type:     "incoming_camera_status",
				FromUser: payload.UserID,
				ChatType: payload.ChatType,
				Enabled:  payload.Enabled,
			})

		case "screen_status":
			var payload common.IncomingScreenStatus
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid screen_status payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			voice.SetCallMediaStatus(c.RDB, roomID, c.UserID, "screen", payload.Enabled)

			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingScreenStatus{
				Type:     "incoming_screen_status",
				FromUser: payload.UserID,
				ChatType: payload.ChatType,
				Enabled:  payload.Enabled,
			})

		case "mic_status":
			var payload common.IncomingMicStatus
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.sendError("invalid mic_status payload")
				continue
			}

			if payload.ChatType != "private" || payload.ReceiverID == "" {
				c.sendError("invalid call context")
				continue
			}

			roomID, ok := c.resolveRoomID(payload.ChatType, payload.ReceiverID)
			if !ok {
				c.sendError("access denied")
				continue
			}

			// сохраняем состояние в Redis (аналогично камере/экрану)
			voice.SetCallMediaStatus(c.RDB, roomID, c.UserID, "mic", payload.Enabled)

			// пересылаем собеседнику
			c.Hub.SendToUser(payload.ReceiverID, common.OutgoingMicStatus{
				Type:     "incoming_mic_status",
				FromUser: payload.UserID,
				ChatType: payload.ChatType,
				Enabled:  payload.Enabled,
			})

		default:
			c.sendError("unsupported message type: " + msgType.Type)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			log.Printf("client rooms: %v", c.Rooms)
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}