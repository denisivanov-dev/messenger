package register

import (
	"messenger/backend_golang/internal/gateway/handlers"
)

func RegisterAllModules() {
	handlers.RegisterInit()
	handlers.RegisterChat()
	handlers.RegisterTyping()
	handlers.RegisterCalls()
	handlers.RegisterWebRTC()
	handlers.RegisterMedia()
}
