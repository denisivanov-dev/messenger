package voice

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func SetCallParticipant(rdb *redis.Client, roomID string, userID string, status string) {
	key := "callroom:" + roomID
	err := rdb.HSet(ctx, key, userID, status).Err()
	if err != nil {
		log.Printf("[voice] redis HSET error: %v", err)
	}
}

func RemoveCallParticipant(rdb *redis.Client, roomID string, userID string) {
	key := "callroom:" + roomID

	if err := rdb.HDel(ctx, key, userID).Err(); err != nil {
		log.Printf("[voice] redis HDEL error: %v", err)
		return
	}

	remaining, err := rdb.HLen(ctx, key).Result()
	if err != nil {
		log.Printf("[voice] redis HLEN error: %v", err)
		return
	}

	if remaining == 0 {
		if err := rdb.Del(ctx, key).Err(); err != nil {
			log.Printf("[voice] redis auto-DEL empty callroom error: %v", err)
		} else {
			log.Printf("[voice] auto-deleted empty callroom: %s", key)
		}
	}
}

func GetCallParticipants(rdb *redis.Client, roomID string) (map[string]string, error) {
	key := "callroom:" + roomID
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("[voice] redis HGETALL error: %v", err)
		return nil, err
	}
	return result, nil
}

func ClearCallRoom(rdb *redis.Client, roomID string) {
	key := "callroom:" + roomID
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("[voice] redis manual DEL error: %v", err)
	}
}