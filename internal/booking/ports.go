package booking

import "context"

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookingsByMovie(movieID string) ([]Booking, error)

	Confirm(ctx context.Context, sessionID, userId string) (Booking, error)
	Release(ctx context.Context, sessionID, userId string) error
}
