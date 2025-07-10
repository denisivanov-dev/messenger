package chat

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
)

func SaveToRedis(rdb *redis.Client, msg common.OutgoingMessage) {
	key := "chat:" + msg.UserID + ":" + msg.ReceiverID

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal redis error: %v", err)
		return
	}

	err = rdb.LPush(context.Background(), key, data).Err() 
	if err != nil {
		log.Printf("redis LPUSH error: %v", err)
	}
}