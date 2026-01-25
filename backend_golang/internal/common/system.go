package common

type SystemPayload struct {
	SystemID   string       `json:"system_id,omitempty"`   // уникальный ID системного сообщения
	Text       string       `json:"text,omitempty"`        // текст системного уведомления
	Event      string       `json:"event,omitempty"`       // "user_joined", "user_left", "call_missed", "server_notice"
	ActorID    string       `json:"actor_id,omitempty"`    // кто вызвал событие (пользователь, система и т.д.)
	TargetID   string       `json:"target_id,omitempty"`   // кого/что касается событие
	Category   string       `json:"category,omitempty"`    // "info" | "warning" | "error" | "event"
	EditedAt   int64        `json:"edited_at,omitempty"`   // если сообщение обновлено
	Pinned     bool         `json:"pinned,omitempty"`      // можно пинить важные уведомления
	Attachments []Attachment `json:"attachments,omitempty"` // файлы/иконки/иллюстрации
	Timestamp  int64        `json:"timestamp,omitempty"`   // время события
	Context    string       `json:"context,omitempty"`     // например chat_id или call_id
}
