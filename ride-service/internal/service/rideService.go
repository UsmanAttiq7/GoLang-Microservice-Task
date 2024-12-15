package service

import (
	"context"
	"fmt"

	pb "github.com/golang_falcon_task/ride-service/proto/ride/v1"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RideService struct {
	db *pgxpool.Pool
	pb.UnimplementedRideServiceServer
}

// NewRideService initializes a new RideService with a database connection.
func NewRideService(db *pgxpool.Pool) *RideService {
	return &RideService{db: db}
}

// UpdateRide updates the details of an existing ride.
func (s *RideService) UpdateRide(ctx context.Context, req *pb.UpdateRideRequest) (*pb.UpdateRideResponse, error) {
	// Input validation
	if req.RideId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ride_id: must be a positive integer")
	}
	if req.Ride == nil {
		return nil, status.Errorf(codes.InvalidArgument, "ride details must be provided")
	}

	// Update the ride details in the database
	_, err := s.db.Exec(ctx, `
        UPDATE rides
        SET source = $1, destination = $2, distance = $3, cost = $4
        WHERE ride_id = $5
    `, req.Ride.Source, req.Ride.Destination, req.Ride.Distance, req.Ride.Cost, req.RideId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "ride with id %d not found", req.RideId)
		}
		return nil, status.Errorf(codes.Internal, "failed to update ride: %v", err)
	}

	// Return a success response
	return &pb.UpdateRideResponse{
		Message: fmt.Sprintf("ride with id %d successfully updated", req.RideId),
	}, nil
}
