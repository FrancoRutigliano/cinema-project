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

func (r *RedisStore) ListBookingsByMovie(movieID string) ([]Booking, error) {
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	var sessions []Booking

	ctx := context.Background()

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := r.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			return nil, fmt.Errorf("getting key %s: %w", iter.Val(), err)
		}
		session, err := parseSession(val)
		if err != nil {
			return nil, fmt.Errorf("parsing session %s: %w", iter.Val(), err)
		}
		sessions = append(sessions, session)
	}

	// error del iterator en sí (ej: Redis se cayó a mitad del scan)
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("scanning keys: %w", err)
	}

	return sessions, nil
}

func (r *RedisStore) Book(book Booking) (Booking, error) {
	session, err := r.hold(book)
	if err != nil {
		return Booking{}, err
	}

	return session, nil
}

func (r *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	session, seatKey, err := r.getSession(ctx, sessionID, userID)
	if err != nil {
		return Booking{}, err
	}

	// Quita el TTL de ambas keys: la seat key y el reverse-lookup de sesión
	// quedan permanentes en Redis (la reserva ya no expira).
	r.rdb.Persist(ctx, seatKey)
	r.rdb.Persist(ctx, sessionKey(sessionID))

	session.Status = StatusConfirmed
	data := Booking{
		ID:      session.ID,
		MovieID: string(session.MovieID),
		SeatID:  session.SeatID,
		UserID:  session.UserID,
		Status:  StatusConfirmed,
	}

	val, err := json.Marshal(&data)
	if err != nil {
		return Booking{}, fmt.Errorf("confirm marshal: %v", err)
	}

	r.rdb.Set(ctx, seatKey, val, 0)

	return session, nil
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

// getSession hace un reverse-lookup: usa el sessionID para obtener la seat key
// (session:<sessionID> → seat:<movieID>:<seatID>), luego lee el valor de esa
// seat key y lo deserializa en un Booking. Devuelve el Booking y la seat key
// para que el caller pueda operar directamente sobre ella (ej: Persist, Del).
func (s *RedisStore) getSession(ctx context.Context, sessionID string, userID string) (Booking, string, error) {
	// Obtiene la seat key asociada al sessionID via el índice inverso.
	seatKey, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return Booking{}, "", fmt.Errorf("getting session key %s: %w", sessionID, err)
	}

	// Lee el valor de la seat key para obtener los datos de la reserva.
	val, err := s.rdb.Get(ctx, seatKey).Result()
	if err != nil {
		return Booking{}, "", fmt.Errorf("getting seat key %s: %w", seatKey, err)
	}

	// Deserializa el JSON almacenado en un Booking.
	session, err := parseSession(val)
	if err != nil {
		return Booking{}, "", fmt.Errorf("parsing session %s: %w", sessionID, err)
	}

	return session, seatKey, nil
}

func parseSession(val string) (Booking, error) {
	var data Booking
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil
}
