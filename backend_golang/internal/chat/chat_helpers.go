package chat

import "messenger/backend_golang/internal/common"

func NormalizeMessage(msg *common.OutgoingMessage, chatID string) {
	if msg.ChatID == "" {
		msg.ChatID = chatID
	}
	if msg.Type == "" {
		if chatID == "1" {
			msg.Type = "global"
		} else {
			msg.Type = "private"
		}
	}
}