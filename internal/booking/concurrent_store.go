package booking

import "sync"

type ConcurrentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: make(map[string]Booking),
	}
}

func (m *ConcurrentStore) Book(book Booking) error {
	if _, exists := m.bookings[book.SeatID]; exists {
		return ErrSeatAlreadyTaken
	}

	m.bookings[book.SeatID] = book

	return nil
}

func (m *ConcurrentStore) ListBookingsByMovie(movieID string) []Booking {
	var result []Booking
	for _, booking := range m.bookings {
		if booking.MovieID == movieID {
			result = append(result, booking)
		}
	}
	return result
}
