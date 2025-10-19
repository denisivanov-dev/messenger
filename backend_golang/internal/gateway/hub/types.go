package hub

const systemRoom = "sys"

type RoomMessage struct {
	RoomID string
	Data   []byte
}

type joinReq struct {
	Client *Client
	RoomID string
}