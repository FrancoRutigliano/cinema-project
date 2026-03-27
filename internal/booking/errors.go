package booking

import "errors"

var (
	ErrSeatAlreadyTaken = errors.New("seat already taken")
	ErrMissingUserID    = errors.New("missing user ID")
)

type ErrorResponse struct {
	Error string `json:"error"`
}
