package common

type Envelope struct {
	Kind      string `json:"kind"`                // "message" | "event" | "call" | "system"
	Action    string `json:"action"`              // "send" | "edit" | "delete" | "typing" | "offer" | ...
	ChatID    string `json:"chat_id,omitempty"`
	ChatType  string `json:"chat_type,omitempty"` // "global" | "private" | "group"
	SenderID  string `json:"sender_id,omitempty"`
	TargetID  string `json:"target_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"` // UNIX time
	Payload   any    `json:"payload,omitempty"`   // (MessagePayload, CallPayload, ...)
}
