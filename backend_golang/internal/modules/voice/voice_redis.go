package voice

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ========== Active call participants ==========

// SetCallParticipant sets a user's participation status ("joined", "calling", etc.)
func SetCallParticipant(rdb *redis.Client, roomID string, userID string, status string) {
	key := "callroom:" + roomID
	err := rdb.HSet(ctx, key, userID, status).Err()
	if err != nil {
		log.Printf("[voice] redis HSET error: %v", err)
	}
}

// RemoveCallParticipant removes a user and their media statuses from the call
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

// GetCallParticipants returns only users with status "joined"
func GetCallParticipants(rdb *redis.Client, roomID string) (map[string]string, error) {
	key := "callroom:" + roomID
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("[voice] redis HGETALL error: %v", err)
		return nil, err
	}

	participants := make(map[string]string)
	for field, val := range result {
		if strings.Contains(field, ":") {
			continue
		}
		if val == "joined" {
			participants[field] = val
		}
	}

	return participants, nil
}

// ClearCallRoom deletes the entire callroom key
func ClearCallRoom(rdb *redis.Client, roomID string) {
	key := "callroom:" + roomID
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("[voice] redis manual DEL error: %v", err)
	}
}

// ========== Media statuses (camera, mic, screen) ==========

// SetCallMediaStatus sets ON/OFF status for a user's media
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

// RemoveCallMediaStatus deletes a specific media status field
func RemoveCallMediaStatus(rdb *redis.Client, roomID string, userID string, mediaType string) {
	key := "callroom:" + roomID
	field := fmt.Sprintf("%s:%s", mediaType, userID)
	if err := rdb.HDel(ctx, key, field).Err(); err != nil {
		log.Printf("[voice] HDEL %s error: %v", field, err)
	}
}

// GetAllCallMediaStatuses returns status maps for cam, mic, and screen
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
		"camera":     cam,
		"microphone": mic,
		"screen":     screen,
	}
}

// ========== Ongoing call message ID ==========

func SetOngoingCallMessageID(rdb *redis.Client, roomID, messageID string) {
	key := fmt.Sprintf("ongoing_call_msg:%s", roomID)
	err := rdb.Set(ctx, key, messageID, 0).Err()
	if err != nil {
		log.Printf("[voice] redis SET ongoing_call_msg error: %v", err)
	}
}

func GetOngoingCallMessageID(rdb *redis.Client, roomID string) (string, error) {
	key := fmt.Sprintf("ongoing_call_msg:%s", roomID)
	return rdb.Get(ctx, key).Result()
}

func ClearOngoingCallMessageID(rdb *redis.Client, roomID string) {
	key := fmt.Sprintf("ongoing_call_msg:%s", roomID)
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("[voice] redis DEL ongoing_call_msg error: %v", err)
	}
}

// ========== Call start time ==========

func SetCallStartTime(rdb *redis.Client, roomID string, timestamp int64) {
	key := fmt.Sprintf("call_start_time:%s", roomID)
	err := rdb.Set(ctx, key, timestamp, 0).Err()
	if err != nil {
		log.Printf("[voice] redis SET call_start_time error: %v", err)
	}
}

func GetCallStartTime(rdb *redis.Client, roomID string) (int64, error) {
	key := fmt.Sprintf("call_start_time:%s", roomID)
	return rdb.Get(ctx, key).Int64()
}

func ClearCallStartTime(rdb *redis.Client, roomID string) {
	key := fmt.Sprintf("call_start_time:%s", roomID)
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("[voice] redis DEL call_start_time error: %v", err)
	}
}

// ========== Call participation history (who ever joined the call) ==========

// AddCallHistoryParticipant adds a user to the historical set of participants
func AddCallHistoryParticipant(rdb *redis.Client, roomID string, userID string) {
	key := fmt.Sprintf("call:was_in_call:%s", roomID)
	err := rdb.SAdd(ctx, key, userID).Err()
	if err != nil {
		log.Printf("[voice] redis SADD call history error: %v", err)
	}
}

// GetCallHistoryParticipants returns all users who have ever joined the call
func GetCallHistoryParticipants(rdb *redis.Client, roomID string) []string {
	key := fmt.Sprintf("call:was_in_call:%s", roomID)
	ids, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		log.Printf("[voice] redis SMEMBERS error: %v", err)
		return []string{}
	}
	return ids
}

// ClearCallHistoryParticipants clears the set of historical participants
func ClearCallHistoryParticipants(rdb *redis.Client, roomID string) {
	key := fmt.Sprintf("call:was_in_call:%s", roomID)
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("[voice] redis DEL call history error: %v", err)
	}
}