package main

import (
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/ws"
	"context"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
        Password: "",
	    DB:       0,
	})

	ctx := context.Background()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Redis недоступен: %v", err)
	}

	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, rdb, w, r)
	})

	err = http.ListenAndServe(":8080", nil)
	// err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}