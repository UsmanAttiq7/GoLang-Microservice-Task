package main

import (
	"github.com/golang_falcon_task/ride-service/internal/logging"
	"github.com/golang_falcon_task/ride-service/internal/metrics"
	"github.com/golang_falcon_task/ride-service/internal/middleware"
	"github.com/golang_falcon_task/ride-service/internal/store"
	"net"

	"github.com/golang_falcon_task/ride-service/internal/config"
	"github.com/golang_falcon_task/ride-service/internal/db"
	"github.com/golang_falcon_task/ride-service/internal/service"
	pb "github.com/golang_falcon_task/ride-service/proto/ride/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	logging.InitLogger()
	log := logging.Logger

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize metrics
	metrics.InitMetrics()
	metrics.StartMetricsServer(":9007")

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	rideStore := store.NewPGRideStore(database)
	rideService := service.NewRideService(rideStore, log)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(log), // Logs all requests and responses
			middleware.MetricsInterceptor(),    // Captures Prometheus metrics
		),
	)
	pb.RegisterRideServiceServer(grpcServer, rideService)

	// Enable reflection for testing
	reflection.Register(grpcServer)

	log.Println("RideService is running on port 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
