package booking

import "time"

const (
	StatusHeld = "held"
)

// Booking representa una reserva de asiento.
type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}
