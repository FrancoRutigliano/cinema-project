package booking

import "errors"

var (
	ErrSeatAlreadyTaken = errors.New("seat already taken")
)

type ErrorResponse struct {
	Error string `json:"error"`
}
