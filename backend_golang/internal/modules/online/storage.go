package online

import (
	"context"
	"fmt"

	rds "github.com/redis/go-redis/v9"
)

// SetStatus sets the user's current online status in Redis
func SetStatus(ctx context.Context, rdb *rds.Client, userID string, status string) error {
	return rdb.HSet(ctx, fmt.Sprintf("user:%s", userID), "status", status).Err()
}

func GetStatus(ctx context.Context, rdb *rds.Client, userID string) (string, error) {
	val, err := rdb.HGet(ctx, fmt.Sprintf("user:%s", userID), "status").Result()
	if err == rds.Nil {
		return "offline", nil
	}
	return val, err
}

// SetOnline marks the user as online.
func SetOnline(ctx context.Context, rdb *rds.Client, userID string) error {
	return SetStatus(ctx, rdb, userID, "online")
}

// SetOffline marks the user as offline.
func SetOffline(ctx context.Context, rdb *rds.Client, userID string) error {
	return SetStatus(ctx, rdb, userID, "offline")
}