package booking

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookingsByMovie(movieID string) ([]Booking, error)
}
