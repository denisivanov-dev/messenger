package redis

import (
	"context"
	"encoding/json"
	"log"

	rds "github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/utils"
)

var ctx = context.Background()

func SaveToRedis(rdb *rds.Client, msg common.OutgoingMessage) {
	var historyKey, queueKey string

	if msg.Type == "private" {
		chatKey := utils.GeneratePrivateChatKey(msg.UserID, msg.ReceiverID)
		historyKey = "chat:history:" + chatKey
		queueKey = "to_save:private:" + chatKey
		msg.ChatID = chatKey
	} else {
		historyKey = "chat:history:1" // global chat id = 1
		queueKey = "to_save:global"
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	rdb.RPush(ctx, historyKey, data)
	rdb.LTrim(ctx, historyKey, -1000, -1)
	rdb.RPush(ctx, queueKey, data)
}

func LoadChatHistory(rdb *rds.Client, chatID string, limit int64) ([]common.OutgoingMessage, error) {
	key := "chat:history:" + chatID

	vals, err := rdb.LRange(ctx, key, -limit, -1).Result()
	if err != nil {
		return nil, err
	}

	msgs := make([]common.OutgoingMessage, 0, len(vals))
	for _, v := range vals {
		var m common.OutgoingMessage
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			log.Printf("redis unmarshal: %v", err)
			continue
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
