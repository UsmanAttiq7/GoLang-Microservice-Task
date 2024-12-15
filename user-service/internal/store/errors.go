package store

import "errors"

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

// ErrDatabaseOperation is returned for generic database operation errors.
var ErrDatabaseOperation = errors.New("database operation failed")
