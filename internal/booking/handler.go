package booking

import (
	"booking-concurrent/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
			Error: fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}

	data := Booking{
		UserID:  req.UserID,
		SeatID:  seatID,
		MovieID: movieID,
	}

	session, err := h.svc.store.Book(data)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, holdResponse{
		SeatID:    seatID,
		MovieID:   session.MovieID,
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
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

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}

	if req.UserID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("invalid request body: %w", ErrMissingUserID),
		})
		return
	}

	session, err := h.svc.store.Confirm(r.Context(), sessionID, req.UserID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, sessionResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    session.SeatID,
		UserID:    req.UserID,
		Status:    session.Status,
	})
}

func (h *handler) ReleaseSession(w http.ResponseWriter, r *http.Request)
