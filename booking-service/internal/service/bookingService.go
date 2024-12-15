package service

import (
	"context"
	"github.com/golang_falcon_task/booking-service/internal/model"
	"github.com/golang_falcon_task/booking-service/internal/store"
	"github.com/sirupsen/logrus"
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
	log          *logrus.Logger
	pb.UnimplementedBookingServiceServer
}

// NewBookingService initializes a new BookingService with a database connection pool.
func NewBookingService(store BookingStore, logger *logrus.Logger) *BookingService {
	return &BookingService{bookingStore: store, log: logger}
}

// CreateBooking handles creating a new booking and its associated ride.
func (s *BookingService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	// Input validation
	if req.UserId <= 0 {
		s.log.Error("Invalid user_id: must be a positive integer", "user_id", req.UserId)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: must be a positive integer")
	}
	if req.Ride == nil {
		s.log.Error("Ride details must be provided", "user_id", req.UserId)
		return nil, status.Errorf(codes.InvalidArgument, "ride details must be provided")
	}

	// Create a new ride
	rideID, err := s.bookingStore.CreateRide(ctx, req.Ride.Source, req.Ride.Destination, req.Ride.Distance, req.Ride.Cost)
	if err != nil {
		s.log.Error("Failed to create ride", "source", req.Ride.Source, "destination", req.Ride.Destination, "error", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create ride: %v", err)
	}

	// Create a new booking
	bookingTime := time.Now()
	bookingId, err := s.bookingStore.CreateBooking(ctx, req.UserId, rideID, bookingTime)
	if err != nil {
		s.log.Error("Failed to create booking", "user_id", req.UserId, "ride_id", rideID, "error", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create booking: %v", err)
	}

	s.log.Info("Booking created successfully", "booking_id", bookingId)

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
		s.log.Error("Invalid booking_id: must be a positive integer", "booking_id", req.BookingId)
		return nil, status.Errorf(codes.InvalidArgument, "invalid booking_id: must be a positive integer")
	}

	// Fetch booking details
	booking, user, ride, err := s.bookingStore.GetBookingDetails(ctx, req.BookingId)
	if err != nil {
		if err == store.ErrBookingNotFound {
			s.log.Error("Booking not found", "booking_id", req.BookingId)
			return nil, status.Errorf(codes.NotFound, "booking with id %d not found", req.BookingId)
		}
		s.log.Error("Failed to fetch booking details", "booking_id", req.BookingId, "error", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to fetch booking: %v", err)
	}

	s.log.Info("Booking details fetched successfully", "booking_id", req.BookingId, "user_id", booking.UserID, "ride_id", booking.RideID)

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
