package store

import "errors"

var (
	ErrRideNotFound      = errors.New("ride not found")
	ErrDatabaseOperation = errors.New("database operation failed")
)
