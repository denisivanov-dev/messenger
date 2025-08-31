package common

type IncomingStartCall struct {
	Type       string `json:"type"`        // "start_call"
	ChatType   string `json:"chat_type"`   // "private"
	ReceiverID string `json:"receiver_id"` // ID собеседника
}

type OutgoingCallNotification struct {
	Type     string `json:"type"`      // "incoming_call"
	FromUser string `json:"from_user"` // ID звонящего
	ChatType string `json:"chat_type"` // "private"
}

type IncomingCancelCall struct {
	Type       string `json:"type"`        // "cancel_call"
	ChatType   string `json:"chat_type"`   // "private"
	ReceiverID string `json:"receiver_id"` // ID собеседника
}

type OutgoingCancelCallNotification struct {
	Type     string `json:"type"`      // "incoming_cancel_call"
	FromUser string `json:"from_user"` // ID звонящего
	ChatType string `json:"chat_type"` // "private"
}

type IncomingCallAnswer struct {
	Type       string `json:"type"`        // "call_answer"
	ChatType   string `json:"chat_type"`   // "private"
	ReceiverID string `json:"receiver_id"` // ID собеседника
	Accepted   bool   `json:"accepted"`    // принял/отклонил
}

type OutgoingCallAnswer struct {
	Type     string `json:"type"`      // "incoming_call_answer"
	FromUser string `json:"from_user"` // ID звонящего
	ChatType string `json:"chat_type"` // "private"
	Accepted bool   `json:"accepted"`  // принял/отклонил
}

type IncomingJoinCall struct {
	Type       string `json:"type"`        // "join_call"
	ChatType   string `json:"chat_type"`   // "private"
	ReceiverID string `json:"receiver_id"` // ID собеседника
}

type OutgoingJoinCallNotification struct {
	Type     string `json:"type"`      // "incoming_join_call"
	FromUser string `json:"from_user"` // ID присоединившегося
	ChatType string `json:"chat_type"` // "private"
}

type IncomingLeaveCall struct {
	Type       string `json:"type"`        // "leave_call"
	ChatType   string `json:"chat_type"`   // "private"
	ReceiverID string `json:"receiver_id"` // ID собеседника
}

type OutgoingLeaveCallNotification struct {
	Type     string `json:"type"`      // "incoming_leave_call"
	FromUser string `json:"from_user"` // ID вышедшего
	ChatType string `json:"chat_type"` // "private"
}

type IncomingWebRTCOffer struct {
	Type       string                 `json:"type"`        // "webrtc_offer"
	ChatType   string                 `json:"chat_type"`   // "private"
	ReceiverID string                 `json:"receiver_id"` // кому
	Offer      map[string]interface{} `json:"offer"`       // SDP offer
}

type OutgoingWebRTCOffer struct {
	Type     string                 `json:"type"`      // "incoming_webrtc_offer"
	FromUser string                 `json:"from_user"` // кто отправил
	ChatType string                 `json:"chat_type"`
	Offer    map[string]interface{} `json:"offer"`
}

type IncomingIceCandidate struct {
	Type       string                 `json:"type"`        // "ice_candidate"
	ChatType   string                 `json:"chat_type"`
	ReceiverID string                 `json:"receiver_id"`
	Candidate  map[string]interface{} `json:"candidate"`
}

type OutgoingIceCandidate struct {
	Type      string                 `json:"type"`      // "incoming_ice_candidate"
	FromUser  string                 `json:"from_user"` // кто отправил
	ChatType  string                 `json:"chat_type"`
	Candidate map[string]interface{} `json:"candidate"`
}

type IncomingWebRTCAnswer struct {
	Type       string                 `json:"type"`        // "webrtc_answer"
	ReceiverID string                 `json:"receiver_id"` // кому
	ChatType   string                 `json:"chat_type"`
	Answer     map[string]interface{} `json:"answer"`      // SDP answer
}

type OutgoingWebRTCAnswer struct {
	Type     string                 `json:"type"`      // "incoming_webrtc_answer"
	FromUser string                 `json:"from_user"` // кто отправил
	ChatType string                 `json:"chat_type"`
	Answer   map[string]interface{} `json:"answer"`    // SDP answer
}

type IncomingCameraStatus struct {
	Type       string `json:"type"`        // "camera_status"
	ChatType   string `json:"chat_type"`   // "private"
	ReceiverID string `json:"receiver_id"` // кому
	UserID     string `json:"user_id"`     // кто отправил
	Enabled    bool   `json:"enabled"`     // включена ли камера
}

type OutgoingCameraStatus struct {
	Type     string `json:"type"`      // "incoming_camera_status"
	FromUser string `json:"from_user"` // кто обновил
	ChatType string `json:"chat_type"` // private
	Enabled  bool   `json:"enabled"`   // true / false
}