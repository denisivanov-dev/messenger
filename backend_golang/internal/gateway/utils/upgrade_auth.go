package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"messenger/backend_golang/internal/jwt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// В проде настроить фильтр по origin!
		return true
	},
}

func UpgradeAndAuth(w http.ResponseWriter, r *http.Request) (*websocket.Conn, string, string, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusUnauthorized)
		return nil, "", "", fmt.Errorf("missing token")
	}

	userID, username, err := jwt.ParseToken(token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, "", "", fmt.Errorf("invalid token")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return nil, "", "", err
	}

	return conn, userID, username, nil
}
