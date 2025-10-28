package common

type FriendRequestPayload struct {
	Type     string `json:"type"`     // "friend_request_sent"
	ToID     string `json:"to_id"`
	FromID   string `json:"from_id"`
}