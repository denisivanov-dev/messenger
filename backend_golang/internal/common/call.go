package common

type CallPayload struct {
	Status       string   `json:"status"`                 // "incoming" | "ongoing" | "ended" | "missed" | "cancelled"
	Participants []string `json:"participants,omitempty"` // users IDs
	StartedAt    int64    `json:"started_at,omitempty"`
	EndedAt      int64    `json:"ended_at,omitempty"`
	Duration     int64    `json:"duration,omitempty"`
}

type WebRTCPayload struct {
	FromUser  string                 `json:"from_user"`
	SDP       map[string]interface{} `json:"sdp,omitempty"`
	Candidate map[string]interface{} `json:"candidate,omitempty"`
}

type DeviceStatePayload struct {
	FromUser string `json:"from_user"`
	Device   string `json:"device"` // "camera" | "mic" | "screen"
	Enabled  bool   `json:"enabled"`
}
