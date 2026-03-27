package booking

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

// Redis store implementa sesiones basadas en reservas respaldadas por redis.
type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb: rdb,
	}
}
