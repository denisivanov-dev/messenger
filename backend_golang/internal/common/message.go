package common

type MessagePayload struct {
	MessageID   string       `json:"message_id"`
	Text        string       `json:"text"`
	ReplyTo     *string      `json:"reply_to"`
	ReplyToText *string      `json:"reply_to_text"`
	ReplyToUser *string      `json:"reply_to_user"`
	EditedAt    int64        `json:"edited_at"`
	IsEdited    bool         `json:"is_edited"`
	Pinned      bool         `json:"pinned"`
	Username    string       `json:"username"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Key          string `json:"key"`
	Type         string `json:"type"` // "image", "file", "audio", ...
	Size         int    `json:"size"`
	OriginalName string `json:"original_name"`
	URL          string `json:"url,omitempty"`
}
