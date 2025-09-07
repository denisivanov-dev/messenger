package voice

import (
	"context"
	"fmt"
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

	RemoveCallMediaStatus(rdb, roomID, userID, "cam")
	RemoveCallMediaStatus(rdb, roomID, userID, "mic")
	RemoveCallMediaStatus(rdb, roomID, userID, "screen")

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

func SetCallMediaStatus(rdb *redis.Client, roomID string, userID string, mediaType string, status bool) {
	key := "callroom:" + roomID
	field := fmt.Sprintf("%s:%s", mediaType, userID)
	value := "off"
	if status {
		value = "on"
	}
	if err := rdb.HSet(ctx, key, field, value).Err(); err != nil {
		log.Printf("[voice] HSET %s error: %v", field, err)
	}
}

func RemoveCallMediaStatus(rdb *redis.Client, roomID string, userID string, mediaType string) {
	key := "callroom:" + roomID
	field := fmt.Sprintf("%s:%s", mediaType, userID)
	if err := rdb.HDel(ctx, key, field).Err(); err != nil {
		log.Printf("[voice] HDEL %s error: %v", field, err)
	}
}

func GetAllCallMediaStatuses(rdb *redis.Client, roomID string) map[string]map[string]bool {
	key := "callroom:" + roomID
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("[voice] HGETALL error: %v", err)
		return nil
	}

	cam := make(map[string]bool)
	mic := make(map[string]bool)
	screen := make(map[string]bool)

	for field, val := range result {
		var mediaType, uid string
		_, err := fmt.Sscanf(field, "%3s:%s", &mediaType, &uid)
		if err != nil {
			continue
		}
		status := val == "on"
		switch mediaType {
		case "cam":
			cam[uid] = status
		case "mic":
			mic[uid] = status
		case "screen":
			screen[uid] = status
		}
	}

	return map[string]map[string]bool{
		"camera": cam,
		"microphone": mic,
		"screen": screen,
	}
}