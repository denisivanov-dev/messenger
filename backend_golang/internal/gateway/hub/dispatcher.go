package hub

import (
	"encoding/json"
	"log"

	"messenger/backend_golang/internal/gateway/client"
)

type HandlerFunc func(c *client.Client, raw json.RawMessage)

// registry — таблица всех зарегистрированных хендлеров.
// Ключ — тип сообщения ("send_message", "start_call", и т.д.).
var registry = make(map[string]HandlerFunc)

// Register — регистрирует новый обработчик по типу.
func Register(msgType string, fn HandlerFunc) {
	if _, exists := registry[msgType]; exists {
		log.Printf("[hub] handler already registered for type: %s", msgType)
		return
	}
	registry[msgType] = fn
	log.Printf("[hub] registered handler: %s", msgType)
}

// Handle — точка входа: получает JSON от клиента, парсит "type" и вызывает нужный хендлер.
func Handle(c *client.Client, raw []byte) {
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