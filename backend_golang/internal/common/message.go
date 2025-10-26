package common

type MessagePayload struct {
	MessageID   string       `json:"message_id,omitempty"`
	Text        string       `json:"text,omitempty"`
	ReplyTo     string       `json:"reply_to,omitempty"`
	ReplyToText string       `json:"reply_to_text,omitempty"`
	ReplyToUser string       `json:"reply_to_user,omitempty"`
	EditedAt    int64        `json:"edited_at,omitempty"`
	Pinned      bool         `json:"pinned,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Key          string `json:"key"`
	Type         string `json:"type"` // "image", "file", "audio", ...
	Size         int    `json:"size"`
	OriginalName string `json:"original_name"`
	URL          string `json:"url,omitempty"`
}