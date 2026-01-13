package typing

import (
	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/types"
)

func HandleTypingEvent(c types.ClientLike, env common.Envelope) {
	if env.ChatType == "" {
		c.SendError("invalid typing envelope")
		return
	}

	switch env.ChatType {

	case "private":
		if env.TargetID == "" {
			c.SendError("invalid typing envelope")
			return
		}

	case "group", "channel":
		if env.ChatID == "" {
			c.SendError("invalid typing envelope")
			return
		}

	case "global":
		// global chat: no target_id / chat_id required

	default:
		c.SendError("invalid typing envelope")
		return
	}
}