package client

import (
	"github.com/gorilla/websocket"
	rds "github.com/redis/go-redis/v9"
)

type RoomMessage struct {
	RoomID string
	Data   []byte
}

type joinReq struct {
	Client *Client
	RoomID string
}

const SystemRoom = "sys"