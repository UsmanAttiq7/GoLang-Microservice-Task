package service

import (
	"context"
	"errors"
	"github.com/golang_falcon_task/ride-service/internal/model"
	"github.com/golang_falcon_task/ride-service/internal/service/mocks"
	"github.com/golang_falcon_task/ride-service/internal/store"
	pb "github.com/golang_falcon_task/ride-service/proto/ride/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestRideService_UpdateRide(t *testing.T) {
	// Create a mock RideStore
	mockStore := new(mocks.RideStore)
	logger := logrus.New()

	// Create a RideService instance
	service := NewRideService(mockStore, logger)

	// Define test cases
	tests := []struct {
		name         string
		rideID       int32
		rideDetails  *pb.Ride
		setupMock    func()
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name:   "Success",
			rideID: 1,
			rideDetails: &pb.Ride{
				Source:      "Downtown",
				Destination: "Airport",
				Distance:    20,
				Cost:        500,
			},
			setupMock: func() {
				mockStore.On("UpdateRide", mock.Anything, int32(1), &model.Ride{
					Source:      "Downtown",
					Destination: "Airport",
					Distance:    20,
					Cost:        500,
				}).Return(nil)
			},
			expectedCode: codes.OK,
			expectedMsg:  "ride with id 1 successfully updated",
		},
		{
			name:         "Invalid RideID",
			rideID:       0,
			rideDetails:  &pb.Ride{},
			setupMock:    func() {},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "",
		},
		{
			name:   "Ride Not Found",
			rideID: 2,
			rideDetails: &pb.Ride{
				Source:      "Downtown",
				Destination: "Airport",
				Distance:    20,
				Cost:        500,
			},
			setupMock: func() {
				mockStore.On("UpdateRide", mock.Anything, int32(2), &model.Ride{
					Source:      "Downtown",
					Destination: "Airport",
					Distance:    20,
					Cost:        500,
				}).Return(store.ErrRideNotFound)
			},
			expectedCode: codes.NotFound,
			expectedMsg:  "",
		},
		{
			name:   "Internal Error",
			rideID: 3,
			rideDetails: &pb.Ride{
				Source:      "Downtown",
				Destination: "Airport",
				Distance:    20,
				Cost:        500,
			},
			setupMock: func() {
				mockStore.On("UpdateRide", mock.Anything, int32(3), &model.Ride{
					Source:      "Downtown",
					Destination: "Airport",
					Distance:    20,
					Cost:        500,
				}).Return(errors.New("database error"))
			},
			expectedCode: codes.Internal,
			expectedMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Call the method
			res, err := service.UpdateRide(context.Background(), &pb.UpdateRideRequest{
				RideId: tt.rideID,
				Ride:   tt.rideDetails,
			})

			// Assert the results
			if tt.expectedCode != codes.OK {
				require.Error(t, err)
				grpcErr, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.expectedCode, grpcErr.Code())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedMsg, res.Message)
			}

			// Assert mock expectations
			mockStore.AssertExpectations(t)
		})
	}
}
