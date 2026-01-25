//go:build ignore

package handlers

import (
	"encoding/json"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/gateway/types"
	gwutils "messenger/backend_golang/internal/gateway/utils"
)

func RegisterWebRTC() {
	hub.Register("webrtc_offer", handleWebRTCOffer)
	hub.Register("webrtc_answer", handleWebRTCAnswer)
	hub.Register("ice_candidate", handleIceCandidate)
}

func handleWebRTCOffer(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingWebRTCOffer
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid webrtc_offer payload") {
		return
	}

	if payload.ReceiverID == "" || payload.Offer == nil {
		c.SendError("invalid webrtc_offer structure")
		return
	}

	c.GetHub().SendToUser(payload.ReceiverID, common.OutgoingWebRTCOffer{
		Type:     "incoming_webrtc_offer",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
		Offer:    payload.Offer,
	})
}

func handleWebRTCAnswer(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingWebRTCAnswer
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid webrtc_answer payload") {
		return
	}

	if payload.ReceiverID == "" || payload.Answer == nil {
		c.SendError("invalid webrtc_answer structure")
		return
	}

	c.GetHub().SendToUser(payload.ReceiverID, common.OutgoingWebRTCAnswer{
		Type:     "incoming_webrtc_answer",
		FromUser: c.ID(),
		ChatType: payload.ChatType,
		Answer:   payload.Answer,
	})
}

func handleIceCandidate(c types.ClientLike, raw json.RawMessage) {
	var payload common.IncomingIceCandidate
	if !gwutils.UnmarshalPayload(raw, &payload, c, "invalid ice_candidate payload") {
		return
	}

	if payload.ReceiverID == "" || payload.Candidate == nil {
		c.SendError("invalid ice_candidate structure")
		return
	}

	c.GetHub().SendToUser(payload.ReceiverID, common.OutgoingIceCandidate{
		Type:      "incoming_ice_candidate",
		FromUser:  c.ID(),
		ChatType:  payload.ChatType,
		Candidate: payload.Candidate,
	})
}
