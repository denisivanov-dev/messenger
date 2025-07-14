package chat

import (
	"encoding/json"
	"log"
	"context"

	rds "github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/utils"
	"messenger/backend_golang/redis"
)

func SaveMessageToCache(rdb *rds.Client, msg common.OutgoingMessage) {
	var chatKey, queueKey string

	if msg.Type == "private" {
		chatKey = utils.GeneratePrivateChatKey(msg.UserID, msg.ReceiverID)
		queueKey = "to_save:private:" + chatKey
	} else {
		chatKey = "1"
		queueKey = "to_save:global"
	}
	msg.ChatID = chatKey

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	historyKey := "chat:history:" + chatKey

	redis.RPush(rdb, historyKey, data)
	redis.LTrim(rdb, historyKey, -1000, -1)
	redis.RPush(rdb, queueKey, data)
}

func DeleteMessageFromRedisHistory(rdb *rds.Client, chatID, messageID string) *common.MessageDeleted {
	key := "chat:history:" + chatID
	queue := "to_delete:global"
	if chatID != "1" {
		queue = "to_delete:private:" + chatID
	}

	vals, err := redis.LRange(rdb, key, 0, -1)
	if err != nil {
		log.Printf("redis lrange error: %v", err)
		return nil
	}

	for _, raw := range vals {
		var msg common.OutgoingMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}
		if msg.MessageID != messageID {
			continue
		}

		if _, err := rdb.LRem(context.Background(), key, 1, raw).Result(); err != nil {
			log.Printf("redis lrem error: %v", err)
		}

		payload := &common.MessageDeleted{
			Type:      "message_deleted",
			MessageID: msg.MessageID,
			ChatID: chatID,
		}

		if data, err := json.Marshal(payload); err == nil {
			redis.RPush(rdb, queue, data)
		}

		return payload
	}

	return nil
}

func LoadMessageHistory(rdb *rds.Client, chatID string, limit int64) ([]common.OutgoingMessage, error) {
	vals, err := redis.LRange(rdb, "chat:history:"+chatID, -limit, -1)
	if err != nil {
		return nil, err
	}

	msgs := make([]common.OutgoingMessage, 0, len(vals))
	for _, v := range vals {
		var m common.OutgoingMessage
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			log.Printf("unmarshal: %v", err)
			continue
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}