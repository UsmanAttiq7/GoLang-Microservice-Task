package store

import (
	"context"
	"fmt"
	"time"

	"github.com/golang_falcon_task/booking-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGBookingStore struct {
	db *pgxpool.Pool
}

// NewPGBookingStore creates a new PGBookingStore instance.
func NewPGBookingStore(db *pgxpool.Pool) *PGBookingStore {
	return &PGBookingStore{db: db}
}

// CreateRide inserts a new ride into the database.
func (s *PGBookingStore) CreateRide(ctx context.Context, source, destination string, distance, cost int32) (int32, error) {
	var rideID int32
	err := s.db.QueryRow(ctx, `
        INSERT INTO rides (source, destination, distance, cost)
        VALUES ($1, $2, $3, $4)
        RETURNING ride_id
    `, source, destination, distance, cost).Scan(&rideID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}
	return rideID, nil
}

// CreateBooking inserts a new booking into the database.
func (s *PGBookingStore) CreateBooking(ctx context.Context, userID, rideID int32, bookingTime time.Time) (int32, error) {
	var bookingID int32
	err := s.db.QueryRow(ctx, `
        INSERT INTO bookings (user_id, ride_id, time)
        VALUES ($1, $2, $3)
        RETURNING booking_id
    `, userID, rideID, bookingTime).Scan(&bookingID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}
	return bookingID, nil
}

// GetBookingDetails retrieves booking details with joins on users and rides tables.
func (s *PGBookingStore) GetBookingDetails(ctx context.Context, bookingID int32) (*model.Booking, *model.User, *model.Ride, error) {
	var (
		booking model.Booking
		user    model.User
		ride    model.Ride
	)

	err := s.db.QueryRow(ctx, `
        SELECT b.booking_id, b.user_id, b.ride_id, b.time, 
               u.user_id, u.name,
               r.ride_id, r.source, r.destination, r.distance, r.cost
        FROM bookings b
        JOIN users u ON b.user_id = u.user_id
        JOIN rides r ON b.ride_id = r.ride_id
        WHERE b.booking_id = $1
    `, bookingID).Scan(
		&booking.ID, &booking.UserID, &booking.RideID, &booking.Timestamp,
		&user.ID, &user.Name,
		&ride.ID, &ride.Source, &ride.Destination, &ride.Distance, &ride.Cost,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil, nil, ErrBookingNotFound
		}
		return nil, nil, nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	return &booking, &user, &ride, nil
}
