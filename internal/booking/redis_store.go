package booking

import (
	"context"
	"encoding/json"
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
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", book.MovieID, book.SeatID)

	book.ID = id
	value, err := json.Marshal(book)
	if err != nil {
		return Booking{}, fmt.Errorf("hold marshaling book: %w", err)
	}

	result := r.rdb.SetArgs(ctx, key, value, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	ok := result.Val() == "OK"
	if !ok {
		return Booking{}, ErrSeatAlreadyTaken
	}

	r.rdb.Set(ctx, sessionKey(key), value, defaultHoldTTL)

	return Booking{
		ID:        id,
		MovieID:   book.MovieID,
		SeatID:    book.SeatID,
		UserID:    book.UserID,
		Status:    StatusHeld,
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}
