package hub

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/common"
	"messenger/backend_golang/internal/gateway/types"
)

type HandlerFunc func(c types.ClientLike, env common.Envelope)

var registry = make(map[string]HandlerFunc)

func Register(key string, fn HandlerFunc) {
	if _, exists := registry[key]; exists {
		log.Printf("[hub] handler already registered for key: %s", key)
		return
	}
	registry[key] = fn
	log.Printf("[hub] registered handler: %s", key)
}

func Handle(c types.ClientLike, raw []byte) {
	var env common.Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		c.SendError("invalid JSON format")
		return
	}

	// --- pretty debug log ---
	if pretty, err := json.MarshalIndent(env, "", "  "); err == nil {
		log.Printf("[gateway] received envelope:\n%s", string(pretty))
	} else {
		log.Printf("[gateway] raw message: %s", string(raw))
	}
	// -------------------------

	if env.Kind == "" || env.Action == "" {
		c.SendError("missing kind or action in envelope")
		return
	}

	// Key example: "message_send", "call_offer", "system_send"
	key := env.Kind + "_" + env.Action

	handler, ok := registry[key]
	if !ok {
		c.SendError("unsupported kind/action: " + key)
		return
	}

	handler(c, env)
}
