package chat

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	rds "github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/redis"
)

func DeleteMessageFromRedisHistory(rdb *rds.Client, chatID, messageID, 
	currentUserID, username string,) *common.Envelope {

	historyKey := HistoryKey(chatID)
	queueKey := QueueKey("to_delete", chatID)

	vals, err := redis.LRange(rdb, historyKey, 0, -1)
	if err != nil {
		log.Printf("[chat][delete] LRange error: %v", err)
		return nil
	}

	for _, raw := range vals {
		var env common.Envelope
		if json.Unmarshal([]byte(raw), &env) != nil {
			continue
		}

		var msg common.MessagePayload
		b, _ := json.Marshal(env.Payload)
		_ = json.Unmarshal(b, &msg)

		if msg.MessageID != messageID || env.SenderID != currentUserID {
			continue
		}

		_, _ = rdb.LRem(context.Background(), historyKey, 1, raw).Result()

		out := BuildMessageEnvelope(
			"delete",
			env.ChatType,
			chatID,
			currentUserID,
			env.TargetID,
			username,
			common.MessagePayload{
				MessageID: messageID,
			},
		)

		if data, err := json.Marshal(out); err == nil {
			redis.RPush(rdb, queueKey, data)
		}
		return &out
	}

	return nil
}

func EditMessageInRedisHistory(rdb *rds.Client,chatID, messageID, 
	newText, currentUserID, username string,) *common.Envelope {

	historyKey := HistoryKey(chatID)
	queueKey := QueueKey("to_edit", chatID)

	vals, err := redis.LRange(rdb, historyKey, 0, -1)
	if err != nil {
		log.Printf("[chat][edit] LRange error: %v", err)
		return nil
	}

	for i, raw := range vals {
		var env common.Envelope
		if json.Unmarshal([]byte(raw), &env) != nil {
			continue
		}

		var msg common.MessagePayload
		b, _ := json.Marshal(env.Payload)
		_ = json.Unmarshal(b, &msg)

		if msg.MessageID != messageID || env.SenderID != currentUserID {
			continue
		}
		if strings.TrimSpace(msg.Text) == strings.TrimSpace(newText) {
			return nil
		}

		msg.Text = newText
		msg.IsEdited = true
		msg.EditedAt = time.Now().UnixMilli()
		env.Payload = msg

		updatedRaw, _ := json.Marshal(env)
		_ = rdb.LSet(context.Background(), historyKey, int64(i), updatedRaw).Err()

		out := BuildMessageEnvelope(
			"edit",
			env.ChatType,
			chatID,
			currentUserID,
			env.TargetID,
			username,
			msg,
		)

		if data, err := json.Marshal(out); err == nil {
			redis.RPush(rdb, queueKey, data)
		}
		return &out
	}

	return nil
}


func PinMessageInRedisHistory(rdb *rds.Client, chatID, messageID string, pin bool, currentUserID string) *common.Envelope {
	historyKey := HistoryKey(chatID)
	queueKey := QueueKey("to_pin", chatID)

	vals, _ := redis.LRange(rdb, historyKey, 0, -1)
	for i, raw := range vals {
		var env common.Envelope
		if err := json.Unmarshal([]byte(raw), &env); err != nil {
			continue
		}

		var msg common.MessagePayload
		b, _ := json.Marshal(env.Payload)
		_ = json.Unmarshal(b, &msg)

		if msg.MessageID != messageID || env.SenderID != currentUserID {
			continue
		}

		msg.Pinned = pin
		env.Payload = msg

		updatedRaw, _ := json.Marshal(env)
		_ = rdb.LSet(context.Background(), historyKey, int64(i), updatedRaw).Err()

		out := common.Envelope{
			Kind:      "message",
			Action:    "pin",
			ChatID:    chatID,
			ChatType:  env.ChatType,
			SenderID:  currentUserID,
			TargetID:  env.TargetID,
			Timestamp: time.Now().UnixMilli(),
			Payload: common.MessagePayload{
				MessageID: msg.MessageID,
				Pinned:    pin,
			},
		}

		if data, err := json.Marshal(out); err == nil {
			redis.RPush(rdb, queueKey, data)
		}
		return &out
	}
	return nil
}