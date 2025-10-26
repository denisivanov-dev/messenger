package chat

import "fmt"

const GlobalRoomID = "1"

func HistoryKey(chatID string) string {
	return "chat:history:" + chatID
}

func QueueKey(kind, chatID string) string {
	if chatID == GlobalRoomID {
		return fmt.Sprintf("%s:global", kind)
	}
	return fmt.Sprintf("%s:private:%s", kind, chatID)
}