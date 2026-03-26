package booking

type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}
