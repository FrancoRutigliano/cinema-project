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

	return nil
}

func (m *MemoryStore) ListBookings(movieID string) []Booking {
	return nil
}
