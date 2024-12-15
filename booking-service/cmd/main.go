package main

import (
	"github.com/golang_falcon_task/booking-service/internal/config"
	"github.com/golang_falcon_task/booking-service/internal/db"
	"github.com/golang_falcon_task/booking-service/internal/service"
	"log"
	"net"

	pb "github.com/golang_falcon_task/booking-service/proto/booking/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, service.NewBookingService(database))

	// Enable reflection for testing
	reflection.Register(grpcServer)

	log.Println("BookingService is running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
