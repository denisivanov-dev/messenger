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

func DeleteMessageFromRedisHistory(rdb *rds.Client, chatID, messageID, currentUserID string) *common.MessageDeleted {
	key := HistoryKey(chatID)
	queue := QueueKey("to_delete", chatID)

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
		if msg.MessageID != messageID || msg.UserID != currentUserID {
			continue
		}

		_, _ = rdb.LRem(context.Background(), key, 1, raw).Result()
		payload := &common.MessageDeleted{
			Type:      "message_deleted",
			MessageID: msg.MessageID,
			ChatID:    chatID,
		}

		if data, err := json.Marshal(payload); err == nil {
			redis.RPush(rdb, queue, data)
		}
		return payload
	}
	return nil
}

func EditMessageInRedisHistory(rdb *rds.Client, chatID, messageID, newText, currentUserID string) *common.MessageEdited {
	historyKey := HistoryKey(chatID)
	queueKey := QueueKey("to_edit", chatID)

	vals, err := redis.LRange(rdb, historyKey, 0, -1)
	if err != nil {
		log.Printf("redis lrange error: %v", err)
		return nil
	}

	for i, raw := range vals {
		var msg common.OutgoingMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}
		if msg.MessageID != messageID || msg.UserID != currentUserID {
			continue
		}

		if strings.TrimSpace(msg.Text) == strings.TrimSpace(newText) {
			return nil
		}

		msg.Text = newText
		msg.EditedAt = time.Now().UnixMilli()
		NormalizeMessage(&msg, chatID)

		updatedRaw, _ := json.Marshal(msg)
		_ = rdb.LSet(context.Background(), historyKey, int64(i), updatedRaw).Err()
		redis.RPush(rdb, queueKey, updatedRaw)

		return &common.MessageEdited{
			Type:      "message_edited",
			MessageID: msg.MessageID,
			ChatID:    chatID,
			NewText:   msg.Text,
			EditedAt:  msg.EditedAt,
		}
	}
	return nil
}

func PinMessageInRedisHistory(rdb *rds.Client, chatID, messageID string, pin bool, currentUserID string) *common.MessagePinned {
	historyKey := HistoryKey(chatID)
	queueKey := QueueKey("to_pin", chatID)

	vals, _ := redis.LRange(rdb, historyKey, 0, -1)
	for i, raw := range vals {
		var msg common.OutgoingMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}
		if msg.MessageID != messageID || msg.UserID != currentUserID {
			continue
		}

		msg.Pinned = pin
		NormalizeMessage(&msg, chatID)
		updatedRaw, _ := json.Marshal(msg)
		_ = rdb.LSet(context.Background(), historyKey, int64(i), updatedRaw).Err()
		redis.RPush(rdb, queueKey, updatedRaw)

		action := "unpin"
		if pin {
			action = "pin"
		}
		return &common.MessagePinned{
			Type:      "message_pinned",
			MessageID: msg.MessageID,
			ChatID:    chatID,
			Action:    action,
		}
	}
	return nil
}

func UpdateSystemMessageInRedisHistory(rdb *rds.Client, chatID, messageID string, updated common.OutgoingMessage) bool {
	historyKey := HistoryKey(chatID)
	queueKey := QueueKey("to_edit", chatID)

	vals, _ := redis.LRange(rdb, historyKey, 0, -1)
	for i, raw := range vals {
		var msg common.OutgoingMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}
		if msg.MessageID != messageID {
			continue
		}

		NormalizeMessage(&updated, chatID)
		updatedRaw, _ := json.Marshal(updated)
		_ = rdb.LSet(context.Background(), historyKey, int64(i), updatedRaw).Err()
		redis.RPush(rdb, queueKey, updatedRaw)
		return true
	}
	return false
}