package common

type FriendPayload struct {
	FromUser string `json:"from_user"`
	ToUser   string `json:"to_user"`
	Status   string `json:"status"` // "pending" | "accepted" | "declined" | "removed"
}
