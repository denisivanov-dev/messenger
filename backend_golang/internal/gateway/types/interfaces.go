package types

import rds "github.com/redis/go-redis/v9"

const SystemRoom = "sys"

type RoomMessage struct {
	RoomID string
	Data   []byte
}

type JoinRequest struct {
	Client ClientLike
	RoomID string
}

type ClientLike interface {
	ID() string
	Username() string
	Send([]byte)
	SendError(string)
	GetRooms() map[string]struct{}

	JoinRoomIfNotJoined(roomID string)
	LeaveAllExcept(allowedRoomIDs ...string)
	BroadcastJSON(roomID string, payload any)

	GetHub() HubLike
	GetRedis() *rds.Client
	GetSendChannel() chan []byte
}

type HubLike interface {
	BroadcastMessage(RoomMessage)
	JoinRoom(JoinRequest)
	UnregisterClient(ClientLike)
	RegisterClient(ClientLike)

	SendToUser(userID string, payload any)
	SendToUsers(userIDs []string, payload any)
}