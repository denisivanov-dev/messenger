package chat

import (
	"fmt"

	"messenger/backend_golang/internal/modules/constants"
)

func HistoryKey(chatID string) string {
	return "chat:history:" + chatID
}

func QueueKey(kind, chatID string) string {
	if chatID == constants.GlobalRoomID {
		return fmt.Sprintf("%s:global", kind)
	}
	return fmt.Sprintf("%s:private:%s", kind, chatID)
}