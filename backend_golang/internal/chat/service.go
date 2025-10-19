package chat

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/utils"
)

// func Parse(raw []byte) (common.IncomingMessage, bool) {
// 	var msg common.IncomingMessage

// 	if err := json.Unmarshal(raw, &msg); err != nil {
// 		return msg, false
// 	}
// 	log.Printf("PARSE DEBUG: %#v", msg)
// 	if strings.TrimSpace(msg.Text) == "" {
// 		return msg, false
// 	}
// 	return msg, true
// }

func BuildMessage(in common.IncomingSendMessage, userID, username string) common.OutgoingMessage {
	t := strings.TrimSpace(strings.ToLower(in.ChatType))

	chatID := ""
	switch t {
	case "global":
		chatID = "1"
	case "private":
		chatID = utils.GeneratePrivateChatKey(userID, in.ReceiverID)
	}

	return common.OutgoingMessage{
		MessageID:    uuid.NewString(),
		ChatID:       chatID,
		Text:         in.Text,
		Attachments:  in.Attachments,
		Timestamp:    time.Now().UnixMilli(),
		Username:     username,
		UserID:       userID,
		Type:         t,
		ReceiverID:   in.ReceiverID,
		ReplyTo:      in.ReplyTo,
		ReplyToText:  in.ReplyToText,
		ReplyToUser:  in.ReplyToUser,
	}
}

func BuildSystemMessage(msgType string, chatType string, fromUserID string, toUserID string, callInfo *common.CallInfo) common.OutgoingMessage {
	t := strings.TrimSpace(strings.ToLower(chatType))

	chatID := ""
	switch t {
	case "global":
		chatID = "1"
	case "private":
		chatID = utils.GeneratePrivateChatKey(fromUserID, toUserID)
	}

	return common.OutgoingMessage{
		MessageID:  uuid.NewString(),
		ChatID:     chatID,
		Type:       msgType,
		UserID:     "0",
		Username:   "system",
		Text:       "",
		Timestamp:  time.Now().UnixMilli(),
		ReceiverID: toUserID,
		CallInfo:   callInfo,
	}
}

package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	rds "github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/utils"
)

// формирует chatKey: private:userA:userB или system-room
func GetRoomKey(userID, chatType, receiverID string) string {
	if chatType == "private" {
		return utils.GeneratePrivateChatKey(userID, receiverID)
	}
	return "1"
}

// проверяет через Redis, имеет ли пользователь доступ к чату
func ResolveRoom(rdb *rds.Client, userID, chatType, receiverID string) (string, bool) {
	ctx := context.Background()
	chatKey := GetRoomKey(userID, chatType, receiverID)

	if chatKey == "1" {
		return chatKey, true
	}

	chatID, err := rdb.Get(ctx, fmt.Sprintf("chat_id:%s", chatKey)).Result()
	if err != nil {
		log.Printf("Redis GET chat_id:%s failed: %v", chatKey, err)
		return "", false
	}

	ok, err := rdb.SIsMember(ctx, fmt.Sprintf("chat:%s:participants", chatID), userID).Result()
	if err != nil || !ok {
		log.Printf("Access denied: user %s not in chat %s", userID, chatID)
		return "", false
	}

	return chatKey, true
}

// отправляет последние сообщения из Redis клиенту
func SendHistory(rdb *rds.Client, roomID string, send chan []byte, limit int) {
	msgs, err := LoadMessageHistory(rdb, roomID, int64(limit))
	if err != nil {
		log.Printf("failed to load chat history for %s: %v", roomID, err)
		return
	}

	for _, msg := range msgs {
		if out, err := json.Marshal(msg); err == nil {
			send <- out
		}
	}
}
