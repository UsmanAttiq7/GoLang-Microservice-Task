package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang_falcon_task/user-service/internal/model"
	"github.com/golang_falcon_task/user-service/internal/store"
	pb "github.com/golang_falcon_task/user-service/proto/user/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserStore defines the interface for user-related operations.
type UserStore interface {
	GetUser(ctx context.Context, id int32) (*model.User, error)
	CreateUser(ctx context.Context, name string) (int32, error)
	DeleteUser(ctx context.Context, id int32) error
}

type UserService struct {
	store UserStore
	log   *logrus.Logger
	pb.UnimplementedUserServiceServer
}

// NewUserService creates a new UserService with a UserStore dependency
func NewUserService(store UserStore, logger *logrus.Logger) *UserService {
	return &UserService{store: store, log: logger}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrUserNotFound):
			s.log.Warn(fmt.Sprintf("user with id %d not found: ", req.UserId), err)
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.UserId)
		default:
			s.log.Error("Failed to get User: ", err)
			return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
		}
	}
	return &pb.GetUserResponse{Name: user.Name}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID, err := s.store.CreateUser(ctx, req.Name)
	if err != nil {
		s.log.Error("Failed to create User: ", err)
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return &pb.CreateUserResponse{UserId: userID}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.store.DeleteUser(ctx, req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrUserNotFound):
			s.log.Warn(fmt.Sprintf("user with id %d not found: ", req.UserId), err)
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.UserId)
		default:
			s.log.Error(fmt.Sprintf("Failed to Delete User with id %d: ", req.UserId), err)
			return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
		}
	}
	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}
