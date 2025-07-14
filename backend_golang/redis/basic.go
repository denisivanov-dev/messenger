package redis

import (
	"context"
	"log"

	rds "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RPush(rdb *rds.Client, key string, data any) {
	err := rdb.RPush(ctx, key, data).Err()
	if err != nil {
		log.Printf("rpush error: %v", err)
	}
}

func LTrim(rdb *rds.Client, key string, start, stop int64) {
	err := rdb.LTrim(ctx, key, start, stop).Err()
	if err != nil {
		log.Printf("ltrim error: %v", err)
	}
}

func LRange(rdb *rds.Client, key string, start, stop int64) ([]string, error) {
	return rdb.LRange(ctx, key, start, stop).Result()
}