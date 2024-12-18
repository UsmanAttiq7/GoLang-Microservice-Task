package main

import (
	"github.com/golang_falcon_task/booking-service/internal/config"
	"github.com/golang_falcon_task/booking-service/internal/db"
	"github.com/golang_falcon_task/booking-service/internal/logging"
	"github.com/golang_falcon_task/booking-service/internal/metrics"
	"github.com/golang_falcon_task/booking-service/internal/middleware"
	"github.com/golang_falcon_task/booking-service/internal/service"
	"github.com/golang_falcon_task/booking-service/internal/store"
	"net"

	pb "github.com/golang_falcon_task/booking-service/proto/booking/v1"
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
	metrics.StartMetricsServer(":9006")

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	bookingStore := store.NewPGBookingStore(database)
	bookingService := service.NewBookingService(bookingStore, log)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(log), // Logs all requests and responses
			middleware.MetricsInterceptor(),    // Captures Prometheus metrics
		),
	)
	pb.RegisterBookingServiceServer(grpcServer, bookingService)

	// Enable reflection for testing
	reflection.Register(grpcServer)

	log.Println("BookingService is running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
