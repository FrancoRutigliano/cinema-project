package main

import (
	"booking-concurrent/internal/adapters/redis"
	"booking-concurrent/internal/booking"
	"booking-concurrent/internal/utils"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", listMovies)

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	store := booking.NewRedisStore(redis.NewClient(
		"localhost:6379"))

	svc := booking.NewService(store)

	handler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/{movieID}/seats", handler.ListSeats)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("shutdown server: %v", err)
	}
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune", Rows: 4, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}
