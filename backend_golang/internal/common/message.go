package common

type IncomingMessage struct {
	Type string `json:"type"`
}

type IncomingSendMessage struct {
	Type        string `json:"type"` // "send_message"
	Text        string `json:"text"`
	ChatType    string `json:"chat_type"`
	ReceiverID  string `json:"receiver_id"`
	ReplyTo     string `json:"reply_to"`
	ReplyToText string `json:"reply_to_text"`
	ReplyToUser string `json:"reply_to_user"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Key          string `json:"key"`
	Type         string `json:"type"` // "image", "file", etc.
	Size         int    `json:"size"`
	OriginalName string `json:"original_name"`
}

type IncomingDeleteMessage struct {
	Type       string `json:"type"` // "delete_message"
	MessageID  string `json:"message_id"`
	ReceiverID string `json:"receiver_id"`
	ChatType   string `json:"chat_type"`
}

type IncomingEditMessage struct {
	Type       string `json:"type"` // "edit_message"
	MessageID  string `json:"message_id"`
	NewText    string `json:"new_text"`
	ReceiverID string `json:"receiver_id"`
	ChatType   string `json:"chat_type"`
}

type IncomingPinMessage struct {
	Type       string `json:"type"` // "pin_message"
	MessageID  string `json:"message_id"`
	ChatID     string `json:"chat_id"`
	Action     string `json:"action"` // "pin" or "unpin"
	ReceiverID string `json:"receiver_id"`
	ChatType   string `json:"chat_type"`
}

type IncomingInitGlobal struct {
	Type     string `json:"type"` // "init_global"
	ChatType string `json:"chat_type"` // "global"
}

type IncomingInitPrivate struct {
	Type       string `json:"type"` // "init_private"
	ChatType   string `json:"chat_type"` // "private"
	ReceiverID string `json:"receiver_id"`
}


type OutgoingMessage struct {
	MessageID   string `json:"message_id"`
	Text        string `json:"text"`
	Timestamp   int64  `json:"timestamp"`
	Username    string `json:"username"`
	UserID      string `json:"user_id"`
	Type        string `json:"type"`
	ReceiverID  string `json:"receiver_id"`
	ChatID      string `json:"chat_id"`
	EditedAt    int64  `json:"edited_at,omitempty"`
	ReplyTo     string `json:"reply_to,omitempty"`
	ReplyToText string `json:"reply_to_text,omitempty"`
	ReplyToUser string `json:"reply_to_user,omitempty"`
	Pinned      bool   `json:"pinned"`
	Attachments   []Attachment `json:"attachments,omitempty"`
}

type TypingMessage struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	ChatID   string `json:"chat_id"`
}

type MessageDeleted struct {
	Type      string `json:"type"`
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
}

type MessageEdited struct {
	Type      string `json:"type"`
	MessageID string `json:"message_id"`
	NewText   string `json:"new_text"`
	EditedAt  int64  `json:"edited_at"`
	ChatID    string `json:"chat_id"`
}

type MessagePinned struct {
	Type      string `json:"type"`
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
	Action    string `json:"action"`
}