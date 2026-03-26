package booking

type MemoryStore struct {
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: make(map[string]Booking),
	}
}

func (m *MemoryStore) Book(book Booking) error {
	if _, exists := m.bookings[book.SeatID]; exists {
		return ErrSeatAlreadyTaken
	}

	m.bookings[book.SeatID] = book

	return nil
}

func (m *MemoryStore) ListBookings(movieID string) []Booking {
	return nil
}
