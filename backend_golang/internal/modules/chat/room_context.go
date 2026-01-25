package chat

import (
	"errors"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/types"
	rds "github.com/redis/go-redis/v9"
)

type RoomContext struct {
	RDB    *rds.Client
	RoomID string
}

func ResolveRoomContext(c types.ClientLike, env common.Envelope) (*RoomContext, error) {
	rdb := c.GetRedis()
	if rdb == nil {
		return nil, errors.New("redis unavailable")
	}

	roomID, ok := ResolveRoom(rdb, c.ID(), env.ChatType, env.TargetID)
	if !ok {
		return nil, errors.New("access denied")
	}

	c.JoinRoomIfNotJoined(roomID)

	return &RoomContext{RDB: rdb, RoomID: roomID}, nil
}
