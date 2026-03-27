package booking

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	defaultHoldTTL = 2 * time.Minute
)

// Redis store implementa sesiones basadas en reservas respaldadas por redis.
type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb: rdb,
	}
}

// Ayuda con el reverse-lookup a buscar sesión por su ID
func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (r *RedisStore) hold(book Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()

	return Booking{
		ID:        id,
		MovieID:   book.MovieID,
		SeatID:    book.SeatID,
		UserID:    book.UserID,
		Status:    StatusHeld,
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}
