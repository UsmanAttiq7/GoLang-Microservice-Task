package service

import (
	"context"
	"github.com/golang_falcon_task/booking-service/internal/model"
	"github.com/golang_falcon_task/booking-service/internal/store"
	"time"

	pb "github.com/golang_falcon_task/booking-service/proto/booking/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingStore interface {
	CreateRide(ctx context.Context, source, destination string, distance, cost int32) (int32, error)
	CreateBooking(ctx context.Context, userID, rideID int32, bookingTime time.Time) (int32, error)
	GetBookingDetails(ctx context.Context, bookingID int32) (*model.Booking, *model.User, *model.Ride, error)
}

type BookingService struct {
	bookingStore BookingStore
	pb.UnimplementedBookingServiceServer
}

// NewBookingService initializes a new BookingService with a database connection pool.
func NewBookingService(store BookingStore) *BookingService {
	return &BookingService{bookingStore: store}
}

// CreateBooking handles creating a new booking and its associated ride.
func (s *BookingService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	// Input validation
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: must be a positive integer")
	}
	if req.Ride == nil {
		return nil, status.Errorf(codes.InvalidArgument, "ride details must be provided")
	}

	// Create a new ride
	rideID, err := s.bookingStore.CreateRide(ctx, req.Ride.Source, req.Ride.Destination, req.Ride.Distance, req.Ride.Cost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create ride: %v", err)
	}

	// Create a new booking
	bookingTime := time.Now()
	bookingId, err := s.bookingStore.CreateBooking(ctx, req.UserId, rideID, bookingTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create booking: %v", err)
	}

	// Return the booking details
	return &pb.CreateBookingResponse{
		Booking: &pb.Booking{
			BookingId: bookingId,
			UserId:    req.UserId,
			RideId:    rideID,
			Time:      bookingTime.Format(time.RFC3339),
		},
	}, nil
}

// GetBooking retrieves booking details, including user and ride information.
func (s *BookingService) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	// Input validation
	if req.BookingId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid booking_id: must be a positive integer")
	}

	// Fetch booking details
	booking, user, ride, err := s.bookingStore.GetBookingDetails(ctx, req.BookingId)
	if err != nil {
		if err == store.ErrBookingNotFound {
			return nil, status.Errorf(codes.NotFound, "booking with id %d not found", req.BookingId)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch booking: %v", err)
	}

	// Return the booking details
	return &pb.GetBookingResponse{
		Name:        user.Name,
		Source:      ride.Source,
		Destination: ride.Destination,
		Distance:    ride.Distance,
		Cost:        ride.Cost,
		Time:        booking.Timestamp.Format(time.RFC3339),
	}, nil
}
