package chat

import (
	"encoding/json"
	"log"

	rds "github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/redis"
)

func SaveMessageToRedisHistory(rdb *rds.Client, roomID string, msg common.Envelope) {
	historyKey := HistoryKey(roomID)
	queueKey := QueueKey("to_save", roomID)

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[chat][SaveMessage] marshal error: %v", err)
		return
	}

	redis.RPush(rdb, historyKey, data)
	redis.LTrim(rdb, historyKey, -1000, -1)

	// push to processing queue to save in DB (handled by Python workers)
	redis.RPush(rdb, queueKey, data)
}

func LoadMessageHistory(rdb *rds.Client, chatID string, limit int64) ([]common.Envelope, error) {
	vals, err := redis.LRange(rdb, HistoryKey(chatID), -limit, -1)
	if err != nil {
		return nil, err
	}

	msgs := make([]common.Envelope, 0, len(vals))
	for _, v := range vals {
		var m common.Envelope
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			log.Printf("[chat][LoadMessage] unmarshal error: %v", err)
			continue
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}