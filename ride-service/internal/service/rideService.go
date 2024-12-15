package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang_falcon_task/ride-service/internal/model"
	"github.com/golang_falcon_task/ride-service/internal/store"
	pb "github.com/golang_falcon_task/ride-service/proto/ride/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RideStore defines the interface for ride-related database operations.
type RideStore interface {
	UpdateRide(ctx context.Context, rideID int32, ride *model.Ride) error
}

type RideService struct {
	rideStore RideStore
	pb.UnimplementedRideServiceServer
}

// NewRideService creates a new RideService with a RideStore dependency.
func NewRideService(store RideStore) *RideService {
	return &RideService{rideStore: store}
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

	// Convert gRPC ride details to model
	ride := &model.Ride{
		Source:      req.Ride.Source,
		Destination: req.Ride.Destination,
		Distance:    req.Ride.Distance,
		Cost:        req.Ride.Cost,
	}

	// Update the ride in the database
	err := s.rideStore.UpdateRide(ctx, req.RideId, ride)
	if err != nil {
		if errors.Is(err, store.ErrRideNotFound) {
			return nil, status.Errorf(codes.NotFound, "ride with id %d not found", req.RideId)
		}
		return nil, status.Errorf(codes.Internal, "failed to update ride: %v", err)
	}

	// Return a success response
	return &pb.UpdateRideResponse{
		Message: fmt.Sprintf("ride with id %d successfully updated", req.RideId),
	}, nil
}
