package common

type OnlinePayload struct {
	UserID string `json:"user_id"`
	Status string `json:"status"` // "online" | "offline" | "away"
}