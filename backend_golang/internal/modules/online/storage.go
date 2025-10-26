package online

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/common"
)

func SetStatus(ctx context.Context, rdb *redis.Client, userID string, status common.Status) error {
	return rdb.HSet(ctx, fmt.Sprintf("user:%s", userID), "status", string(status)).Err()
}