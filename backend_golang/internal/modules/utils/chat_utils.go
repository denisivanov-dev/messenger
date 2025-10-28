package utils

import (
	"messenger/backend_golang/internal/modules/constants"
)

func ResolveChatID(chatType, fromID, toID string) string {
	switch chatType {
	case "global":
		return constants.GlobalRoomID
	case "private":
		return GeneratePrivateChatKey(fromID, toID)
	default:
		return constants.GlobalRoomID
	}
}
