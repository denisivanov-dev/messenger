package chat

import (
	"encoding/json"

	"messenger/backend_golang/redis"
	rds "github.com/redis/go-redis/v9"
)

func LoadHistoryJSON(rdb *rds.Client, chatID string, limit int) ([][]byte, error) {
	msgs, err := redis.LoadChatHistory(rdb, chatID, int64(limit))
	if err != nil {
		return nil, err
	}

	out := make([][]byte, 0, len(msgs))
	for _, m := range msgs {
		b, err := json.Marshal(m)
		if err != nil {
			continue
		}
		out = append(out, b)
	}
	
	return out, nil
} 