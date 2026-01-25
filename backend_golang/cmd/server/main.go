package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"messenger/backend_golang/internal/gateway/handlers/register"
	"messenger/backend_golang/internal/gateway/hub"
	"messenger/backend_golang/internal/pubsub"
	"messenger/backend_golang/internal/gateway/ws"
)

func main() {
	_ = godotenv.Load("../.env")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "0.0.0.0:8080"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis unavailable: %v", err)
	}

	register.RegisterAllModules()

	h := hub.NewHub()

	pubsub.StartFriendRequestSubscriber(ctx, rdb, h)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(h, rdb, w, r)
	})

	log.Printf("Server listening on %s...", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
