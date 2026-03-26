package booking

type BookingStore interface {
	Book(b Booking) error
	ListBookingsByMovie(movieID string) []Booking
}
