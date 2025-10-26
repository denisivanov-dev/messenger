package register

import (
	"log"

	"messenger/backend_golang/internal/gateway/handlers"
)

func RegisterAllModules() {
	log.Println("[register] initializing all WebSocket handlers...")

	handlers.RegisterInit()     // init_global / init_private
	handlers.RegisterChat()     // send_message / edit / delete / pin
	handlers.RegisterTyping()   // typing events
	handlers.RegisterCalls()    // start_call / join / leave / cancel / answer
	handlers.RegisterWebRTC()   // webrtc_offer / answer / ice_candidate
	handlers.RegisterMedia()    // camera / screen / mic status updates

	log.Println("[register] all WebSocket modules registered successfully")
}
