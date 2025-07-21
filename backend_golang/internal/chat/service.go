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