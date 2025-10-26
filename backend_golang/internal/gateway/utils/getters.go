package utils

import (
	rds "github.com/redis/go-redis/v9"
	"messenger/backend_golang/internal/gateway/types"
)

// GetHub safely extracts the Hub from a ClientLike if available.
func GetHub(c types.ClientLike) types.HubLike {
	type hubGetter interface {
		GetHub() types.HubLike
	}
	if h, ok := any(c).(hubGetter); ok {
		return h.GetHub()
	}
	return nil
}

// GetRedis safely extracts *redis.Client from the underlying client.
func GetRedis(c types.ClientLike) *rds.Client {
	type redisGetter interface {
		GetRedis() *rds.Client
	}
	if r, ok := any(c).(redisGetter); ok {
		return r.GetRedis()
	}
	return nil
}

// GetSendChan extracts the Send channel from a ClientLike if available.
func GetSendChan(c types.ClientLike) chan []byte {
	type sender interface {
		GetSendChannel() chan []byte
	}
	if s, ok := any(c).(sender); ok {
		return s.GetSendChannel()
	}
	return nil
}