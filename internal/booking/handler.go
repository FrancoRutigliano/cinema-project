package booking

import (
	"booking-concurrent/internal/utils"
	"encoding/json"
	"fmt"
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

func (h *handler) HoldSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("invalid request body: %w", err),
		})
		return
	}

	h.svc.store.Book(Booking{})
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")

	bookings, err := h.svc.store.ListBookingsByMovie(movieID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
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
