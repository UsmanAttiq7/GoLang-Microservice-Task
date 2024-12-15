package service

import (
	"context"
	pb "github.com/golang_falcon_task/user-service/proto/user/v1"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	db *pgxpool.Pool
	pb.UnimplementedUserServiceServer
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var name string

	// Validate
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be a positive integer")
	}

	err := s.db.QueryRow(ctx, "SELECT name FROM users WHERE user_id = $1", req.UserId).Scan(&name)
	if err == pgx.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.UserId)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}
	// Return the user's name
	return &pb.GetUserResponse{Name: name}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name cannot be empty")
	}
	var userID int32
	err := s.db.QueryRow(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING user_id", req.Name).Scan(&userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUserResponse{UserId: userID}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Validate
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be a positive integer")
	}

	result, err := s.db.Exec(ctx, "DELETE FROM users WHERE user_id = $1", req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.UserId)
	}

	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}
