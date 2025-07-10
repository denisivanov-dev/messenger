package common

type Status string

const (
	Online  Status = "online"
	Offline Status = "offline"
)

type StatusMessage struct {
	Event  string `json:"event"`
	UserID string `json:"user_id"`
	Status Status `json:"status"`
}