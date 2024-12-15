package main

import (
	"github.com/golang_falcon_task/user-service/internal/config"
	"github.com/golang_falcon_task/user-service/internal/db"
	"github.com/golang_falcon_task/user-service/internal/service"
	"github.com/golang_falcon_task/user-service/internal/store"
	pb "github.com/golang_falcon_task/user-service/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userStore := store.NewPGUserStore(database)
	pb.RegisterUserServiceServer(grpcServer, service.NewUserService(userStore))

	// Enable server reflection
	reflection.Register(grpcServer)

	log.Println("UserService is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
