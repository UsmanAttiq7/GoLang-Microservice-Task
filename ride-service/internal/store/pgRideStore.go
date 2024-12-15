package store

import (
	"context"
	"fmt"

	"github.com/golang_falcon_task/ride-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRideStore struct {
	db *pgxpool.Pool
}

// NewPGRideStore creates a new PGRideStore instance.
func NewPGRideStore(db *pgxpool.Pool) *PGRideStore {
	return &PGRideStore{db: db}
}

// UpdateRide updates the details of an existing ride.
func (s *PGRideStore) UpdateRide(ctx context.Context, rideID int32, ride *model.Ride) error {
	result, err := s.db.Exec(ctx, `
        UPDATE rides
        SET source = $1, destination = $2, distance = $3, cost = $4
        WHERE ride_id = $5
    `, ride.Source, ride.Destination, ride.Distance, ride.Cost, rideID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if result.RowsAffected() == 0 {
		return ErrRideNotFound
	}

	return nil
}
