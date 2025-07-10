package online

import (
	"context"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
)

func SetOnline(ctx context.Context, rdb *redis.Client, userID string) error {
	return SetStatus(ctx, rdb, userID, common.Online)
}

func SetOffline(ctx context.Context, rdb *redis.Client, userID string) error {
	return SetStatus(ctx, rdb, userID, common.Offline)
}