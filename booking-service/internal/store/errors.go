package store

import "errors"

// ErrBookingNotFound is returned when a booking is not found.
var ErrBookingNotFound = errors.New("booking not found")

// ErrDatabaseOperation is returned for generic database operation errors.
var ErrDatabaseOperation = errors.New("database operation failed")
