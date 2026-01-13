package common

type OnlinePayload struct {
	Status string `json:"status"` // "online" | "offline" | "away"
}