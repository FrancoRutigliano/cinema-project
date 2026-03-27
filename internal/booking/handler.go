package booking

import (
	"booking-concurrent/internal/utils"
	"net/http"
)

type handler struct {
	svc Service
}

func NewHandler(svc *Service) *handler {
	return &handler{
		svc: *svc,
	}
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")

	bookings, err := h.svc.store.ListBookingsByMovie(movieID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, nil)
	}

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID: b.SeatID,
			UserID: b.UserID,
			Booked: true,
		})
	}

	utils.WriteJSON(w, http.StatusOK, seats)
}
