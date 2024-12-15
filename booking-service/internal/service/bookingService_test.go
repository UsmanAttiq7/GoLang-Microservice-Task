package service

import (
	"context"
	"errors"
	"github.com/golang_falcon_task/booking-service/internal/model"
	"github.com/golang_falcon_task/booking-service/internal/service/mocks"
	"github.com/golang_falcon_task/booking-service/internal/store"
	pb "github.com/golang_falcon_task/booking-service/proto/booking/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestBookingService_CreateBooking(t *testing.T) {
	logger := logrus.New()

	tests := []struct {
		name         string
		req          *pb.CreateBookingRequest
		setupMock    func(mockStore *mocks.BookingStore)
		expectedCode codes.Code
		expectedResp *pb.CreateBookingResponse
	}{
		{
			name: "Success",
			req: &pb.CreateBookingRequest{
				UserId: 1,
				Ride: &pb.Ride{
					Source:      "Downtown",
					Destination: "Airport",
					Distance:    20,
					Cost:        500,
				},
			},
			setupMock: func(mockStore *mocks.BookingStore) {
				mockStore.On("CreateRide", mock.Anything, "Downtown", "Airport", int32(20), int32(500)).
					Return(int32(101), nil)
				mockStore.On("CreateBooking", mock.Anything, int32(1), int32(101), mock.Anything).
					Return(int32(1001), nil)
			},
			expectedCode: codes.OK,
			expectedResp: &pb.CreateBookingResponse{
				Booking: &pb.Booking{
					BookingId: 1001,
					UserId:    1,
					RideId:    101,
					Time:      time.Now().Format(time.RFC3339), // Mock dynamic time if needed
				},
			},
		},
		{
			name: "Ride Creation Failure",
			req: &pb.CreateBookingRequest{
				UserId: 1,
				Ride: &pb.Ride{
					Source:      "Downtown",
					Destination: "Airport",
					Distance:    20,
					Cost:        500,
				},
			},
			setupMock: func(mockStore *mocks.BookingStore) {
				mockStore.On("CreateRide", mock.Anything, "Downtown", "Airport", int32(20), int32(500)).
					Return(int32(0), errors.New("ride creation failed"))
			},
			expectedCode: codes.Internal,
			expectedResp: nil,
		},
		{
			name: "Booking Creation Failure",
			req: &pb.CreateBookingRequest{
				UserId: 1,
				Ride: &pb.Ride{
					Source:      "Downtown",
					Destination: "Airport",
					Distance:    20,
					Cost:        500,
				},
			},
			setupMock: func(mockStore *mocks.BookingStore) {
				mockStore.On("CreateRide", mock.Anything, "Downtown", "Airport", int32(20), int32(500)).
					Return(int32(101), nil)
				mockStore.On("CreateBooking", mock.Anything, int32(1), int32(101), mock.Anything).
					Return(int32(0), errors.New("booking creation failed"))
			},
			expectedCode: codes.Internal,
			expectedResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mockStore for each test case
			mockStore := new(mocks.BookingStore)
			tt.setupMock(mockStore)

			// Create a new service for each test case
			service := NewBookingService(mockStore, logger)

			// Call the method
			resp, err := service.CreateBooking(context.Background(), tt.req)

			// Validate the response
			if tt.expectedCode != codes.OK {
				require.Error(t, err, "Expected an error but got none")
				grpcErr, ok := status.FromError(err)
				require.True(t, ok, "Expected gRPC error")
				require.Equal(t, tt.expectedCode, grpcErr.Code(), "Expected gRPC code mismatch")
			} else {
				require.NoError(t, err, "Expected no error but got one")
				require.NotNil(t, resp, "Expected a valid response but got nil")
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestBookingService_GetBooking(t *testing.T) {
	mockStore := new(mocks.BookingStore)
	logger := logrus.New()
	service := NewBookingService(mockStore, logger)

	tests := []struct {
		name         string
		bookingID    int32
		setupMock    func()
		expectedCode codes.Code
		expectedResp *pb.GetBookingResponse
	}{
		{
			name:      "Success",
			bookingID: 1,
			setupMock: func() {
				mockStore.On("GetBookingDetails", mock.Anything, int32(1)).Return(
					&model.Booking{ID: 1, UserID: 10, RideID: 20, Timestamp: time.Now()},
					&model.User{Name: "John Doe"},
					&model.Ride{Source: "Downtown", Destination: "Airport", Distance: 20, Cost: 500},
					nil,
				)
			},
			expectedCode: codes.OK,
			expectedResp: &pb.GetBookingResponse{
				Name:        "John Doe",
				Source:      "Downtown",
				Destination: "Airport",
				Distance:    20,
				Cost:        500,
				Time:        time.Now().Format(time.RFC3339), // This would vary
			},
		},
		{
			name:         "Invalid Booking ID",
			bookingID:    0,
			setupMock:    func() {},
			expectedCode: codes.InvalidArgument,
			expectedResp: nil,
		},
		{
			name:      "Booking Not Found",
			bookingID: 2,
			setupMock: func() {
				mockStore.On("GetBookingDetails", mock.Anything, int32(2)).Return(nil, nil, nil, store.ErrBookingNotFound)
			},
			expectedCode: codes.NotFound,
			expectedResp: nil,
		},
		{
			name:      "Internal Error",
			bookingID: 3,
			setupMock: func() {
				mockStore.On("GetBookingDetails", mock.Anything, int32(3)).Return(nil, nil, nil, errors.New("database error"))
			},
			expectedCode: codes.Internal,
			expectedResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := service.GetBooking(context.Background(), &pb.GetBookingRequest{BookingId: tt.bookingID})

			if tt.expectedCode != codes.OK {
				require.Error(t, err)
				grpcErr, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.expectedCode, grpcErr.Code())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

			mockStore.AssertExpectations(t)
		})
	}
}
