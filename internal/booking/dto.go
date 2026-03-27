package booking

type seatInfo struct {
	SeatID string `json:"seat_id"`
	UserID string `json:"user_id"`
	Booked bool   `json:"booked"`
}
