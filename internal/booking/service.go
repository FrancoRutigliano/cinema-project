package booking

import "context"

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Book(book Booking) (Booking, error) {
	return s.store.Book(book)
}

func (s *Service) ListBookingsByMovie(movieID string) ([]Booking, error) {
	return s.store.ListBookingsByMovie(movieID)
}

func (s *Service) ConfirmSeat(ctx context.Context, sessionID string, userID string) (Booking, error) {
	return s.store.Confirm(ctx, sessionID, userID)
}
