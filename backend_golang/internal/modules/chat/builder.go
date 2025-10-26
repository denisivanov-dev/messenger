package chat

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/modules/utils"
)

func BuildMessage(in common.IncomingSendMessage, userID, username string) common.Envelope {
	t := strings.TrimSpace(strings.ToLower(in.ChatType))

	var chatID string
	switch t {
	case "global":
		chatID = GlobalRoomID
	case "private":
		chatID = utils.GeneratePrivateChatKey(userID, in.ReceiverID)
	default:
		chatID = GlobalRoomID
	}

	msg := common.MessagePayload{
		MessageID:   uuid.NewString(),
		Text:        in.Text,
		ReplyTo:     in.ReplyTo,
		ReplyToText: in.ReplyToText,
		ReplyToUser: in.ReplyToUser,
		Attachments: in.Attachments,
	}

	return common.Envelope{
		Kind:      "message",
		Action:    "send",
		ChatID:    chatID,
		ChatType:  t,
		SenderID:  userID,
		TargetID:  in.ReceiverID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   msg,
	}
}

func BuildSystemMessage(action, chatType, fromUserID, toUserID string, callInfo *common.CallInfo) common.Envelope {
	t := strings.TrimSpace(strings.ToLower(chatType))

	var chatID string
	switch t {
	case "global":
		chatID = GlobalRoomID
	case "private":
		chatID = utils.GeneratePrivateChatKey(fromUserID, toUserID)
	default:
		chatID = GlobalRoomID
	}

	sysMsg := common.MessagePayload{
		Text: "",
	}

	return common.Envelope{
		Kind:      "system",
		Action:    action, // "call_started", "call_ended", "user_joined"
		ChatID:    chatID,
		ChatType:  t,
		SenderID:  "0",
		TargetID:  toUserID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   sysMsg,
	}
}
