package chat

import (
	"encoding/json"
	"log"
	"context"
	"time"

	rds "github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/redis"
)

func SaveMessageToRedisHistory(rdb *rds.Client, roomID string, msg common.OutgoingMessage) {
	msg.ChatID = roomID

	historyKey := "chat:history:" + roomID
	queueKey := "to_save:global"
	if roomID != "1" {
		queueKey = "to_save:private:" + roomID
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	redis.RPush(rdb, historyKey, data)
	redis.LTrim(rdb, historyKey, -1000, -1)
	redis.RPush(rdb, queueKey, data)
}

func DeleteMessageFromRedisHistory(rdb *rds.Client, chatID, messageID, currentUserID string) *common.MessageDeleted {
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

		if msg.UserID != currentUserID {
			log.Printf("unauthorized deletion attempt by %s for message %s", currentUserID, messageID)
			return nil
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

func EditMessageInRedisHistory(rdb *rds.Client, chatID, messageID, newText, currentUserID string) *common.MessageEdited {
	historyKey := "chat:history:" + chatID
	queueKey := "to_edit:global"
	if chatID != "1" {
		queueKey = "to_edit:private:" + chatID
	}

	log.Printf("EditMessageInRedisHistory: chatID=%s, messageID=%s, newText=%q, currentUserID=%s",
		chatID, messageID, newText, currentUserID)
		
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
		if msg.MessageID != messageID {
			continue
		}

		if msg.UserID != currentUserID {
			log.Printf("unauthorized edit attempt by %s for message %s", currentUserID, messageID)
			return nil
		}

		msg.Text = newText
		msg.EditedAt = time.Now().UnixMilli()

		updatedRaw, err := json.Marshal(msg)
		if err != nil {
			log.Printf("marshal edited message error: %v", err)
			return nil
		}

		if err := rdb.LSet(context.Background(), historyKey, int64(i), updatedRaw).Err(); err != nil {
			log.Printf("redis lset error: %v", err)
			return nil
		}

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