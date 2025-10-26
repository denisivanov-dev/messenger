package chat

import "messenger/backend_golang/internal/common"

// NormalizeMessage ensures that ChatID and Type are filled correctly
func NormalizeMessage(msg *common.OutgoingMessage, chatID string) {
	if msg.ChatID == "" {
		msg.ChatID = chatID
	}
	if msg.Type == "" {
		if chatID == GlobalRoomID {
			msg.Type = "global"
		} else {
			msg.Type = "private"
		}
	}
}