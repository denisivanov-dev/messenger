package common

type IncomingMessage struct {
	Type       string `json:"type"`
	ReceiverID string `json:"receiver_id"`
	ChatType   string `json:"chat_type"`
	Text       string `json:"text"`
	NewText    string `json:"new_text"`
	Timestamp  int64  `json:"timestamp"`
	MessageID  string `json:"message_id"`
}  

type OutgoingMessage struct {
	MessageID  string `json:"message_id"`
	Text       string `json:"text"`
	Timestamp  int64  `json:"timestamp"`
	Username   string `json:"username"`
	UserID     string `json:"user_id"`
	Type       string `json:"type"`      
	ReceiverID string `json:"receiver_id"` 
	ChatID     string `json:"chat_id"`   
	EditedAt   int64  `json:"edited_at,omitempty"`  
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