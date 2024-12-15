package service

import (
	"context"
	"errors"
	"github.com/golang_falcon_task/user-service/internal/model"
	"github.com/golang_falcon_task/user-service/internal/service/mocks"
	"github.com/golang_falcon_task/user-service/internal/store"
	pb "github.com/golang_falcon_task/user-service/proto/user/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestUserService_GetUser(t *testing.T) {
	// Create a mock UserStore
	mockStore := new(mocks.UserStore)
	logger := logrus.New()

	// Create a UserService instance
	service := NewUserService(mockStore, logger)

	// Define test cases
	tests := []struct {
		name         string
		userID       int32
		setupMock    func()
		expectedCode codes.Code
		expectedName string
	}{
		{
			name:   "Success",
			userID: 1,
			setupMock: func() {
				mockStore.On("GetUser", mock.Anything, int32(1)).Return(&model.User{ID: 1, Name: "John Doe"}, nil)
			},
			expectedCode: codes.OK,
			expectedName: "John Doe",
		},
		{
			name:   "User Not Found",
			userID: 2,
			setupMock: func() {
				mockStore.On("GetUser", mock.Anything, int32(2)).Return(nil, store.ErrUserNotFound)
			},
			expectedCode: codes.NotFound,
			expectedName: "",
		},
		{
			name:   "Internal Error",
			userID: 3,
			setupMock: func() {
				mockStore.On("GetUser", mock.Anything, int32(3)).Return(nil, errors.New("database error"))
			},
			expectedCode: codes.Internal,
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Call the method
			res, err := service.GetUser(context.Background(), &pb.GetUserRequest{UserId: tt.userID})

			// Assert the results
			if tt.expectedCode != codes.OK {
				require.Error(t, err)
				grpcErr, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.expectedCode, grpcErr.Code())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedName, res.Name)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	mockStore := new(mocks.UserStore)
	logger := logrus.New()
	service := NewUserService(mockStore, logger)

	tests := []struct {
		name         string
		userName     string
		setupMock    func()
		expectedCode codes.Code
		expectedID   int32
	}{
		{
			name:     "Success",
			userName: "John Doe",
			setupMock: func() {
				mockStore.On("CreateUser", mock.Anything, "John Doe").Return(int32(1), nil)
			},
			expectedCode: codes.OK,
			expectedID:   1,
		},
		{
			name:     "Internal Error",
			userName: "Jane Doe",
			setupMock: func() {
				mockStore.On("CreateUser", mock.Anything, "Jane Doe").Return(int32(0), errors.New("database error"))
			},
			expectedCode: codes.Internal,
			expectedID:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			res, err := service.CreateUser(context.Background(), &pb.CreateUserRequest{Name: tt.userName})

			if tt.expectedCode != codes.OK {
				require.Error(t, err)
				grpcErr, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.expectedCode, grpcErr.Code())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, res.UserId)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	mockStore := new(mocks.UserStore)
	logger := logrus.New()
	service := NewUserService(mockStore, logger)

	tests := []struct {
		name         string
		userID       int32
		setupMock    func()
		expectedCode codes.Code
	}{
		{
			name:   "Success",
			userID: 1,
			setupMock: func() {
				mockStore.On("DeleteUser", mock.Anything, int32(1)).Return(nil)
			},
			expectedCode: codes.OK,
		},
		{
			name:   "User Not Found",
			userID: 2,
			setupMock: func() {
				mockStore.On("DeleteUser", mock.Anything, int32(2)).Return(store.ErrUserNotFound)
			},
			expectedCode: codes.NotFound,
		},
		{
			name:   "Internal Error",
			userID: 3,
			setupMock: func() {
				mockStore.On("DeleteUser", mock.Anything, int32(3)).Return(errors.New("database error"))
			},
			expectedCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			_, err := service.DeleteUser(context.Background(), &pb.DeleteUserRequest{UserId: tt.userID})

			if tt.expectedCode != codes.OK {
				require.Error(t, err)
				grpcErr, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.expectedCode, grpcErr.Code())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
