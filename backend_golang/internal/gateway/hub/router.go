package hub

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/gateway/types"
)

type HandlerFunc func(c types.ClientLike, raw json.RawMessage)

var registry = make(map[string]HandlerFunc)

func Register(msgType string, fn HandlerFunc) {
	if _, exists := registry[msgType]; exists {
		log.Printf("[hub] handler already registered for type: %s", msgType)
		return
	}
	registry[msgType] = fn
	log.Printf("[hub] registered handler: %s", msgType)
}

func Handle(c types.ClientLike, raw []byte) {
	var base struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(raw, &base); err != nil {
		c.SendError("invalid JSON format")
		return
	}

	handler, ok := registry[base.Type]
	if !ok {
		c.SendError("unsupported message type: " + base.Type)
		return
	}

	handler(c, json.RawMessage(raw))
}