package service

import (
	"context"
	"time"

	pb "github.com/golang_falcon_task/booking-service/proto/booking/v1"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingService struct {
	db *pgxpool.Pool
	pb.UnimplementedBookingServiceServer
}

// NewBookingService initializes a new BookingService with a database connection pool.
func NewBookingService(db *pgxpool.Pool) *BookingService {
	return &BookingService{db: db}
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

	// Insert a new ride into the database
	var rideID int32
	err := s.db.QueryRow(ctx, `
        INSERT INTO rides (source, destination, distance, cost)
        VALUES ($1, $2, $3, $4)
        RETURNING ride_id
    `, req.Ride.Source, req.Ride.Destination, req.Ride.Distance, req.Ride.Cost).Scan(&rideID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create ride: %v", err)
	}

	// Insert the booking into the database
	var bookingID int32
	currentTime := time.Now().Format(time.RFC3339)
	err = s.db.QueryRow(ctx, `
        INSERT INTO bookings (user_id, ride_id, time)
        VALUES ($1, $2, $3)
        RETURNING booking_id
    `, req.UserId, rideID, currentTime).Scan(&bookingID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create booking: %v", err)
	}

	// Return the newly created booking details
	return &pb.CreateBookingResponse{
		Booking: &pb.Booking{
			UserId: req.UserId,
			RideId: rideID,
			Time:   currentTime,
		},
	}, nil
}

// GetBooking retrieves booking details, including user and ride information.
func (s *BookingService) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	// Input validation
	if req.BookingId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid booking_id: must be a positive integer")
	}

	var (
		name, source, destination string
		distance, cost            int
		bookingTime               time.Time // Use time.Time for scanning timestamp
	)

	// Query the database for the booking details
	err := s.db.QueryRow(ctx, `
        SELECT u.name, r.source, r.destination, r.distance, r.cost, b.time
        FROM bookings b
        JOIN users u ON b.user_id = u.user_id
        JOIN rides r ON b.ride_id = r.ride_id
        WHERE b.booking_id = $1
    `, req.BookingId).Scan(&name, &source, &destination, &distance, &cost, &bookingTime)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "booking with id %d not found", req.BookingId)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch booking: %v", err)
	}

	// Format the time as RFC3339
	formattedTime := bookingTime.Format(time.RFC3339)

	// Return the booking details
	return &pb.GetBookingResponse{
		Name:        name,
		Source:      source,
		Destination: destination,
		Distance:    int32(distance),
		Cost:        int32(cost),
		Time:        formattedTime,
	}, nil
}
