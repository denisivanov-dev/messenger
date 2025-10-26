package chat

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/modules/utils"
)

func BuildMessage(in common.IncomingSendMessage, userID, username string) common.OutgoingMessage {
	t := strings.TrimSpace(strings.ToLower(in.ChatType))

	chatID := ""
	switch t {
	case "global":
		chatID = GlobalRoomID
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

// BuildSystemMessage generates a system message (for example, when call starts)
func BuildSystemMessage(msgType, chatType, fromUserID, toUserID string, callInfo *common.CallInfo) common.OutgoingMessage {
	t := strings.TrimSpace(strings.ToLower(chatType))

	chatID := ""
	switch t {
	case "global":
		chatID = GlobalRoomID
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