package main

import (
	"github.com/golang_falcon_task/user-service/internal/config"
	"github.com/golang_falcon_task/user-service/internal/db"
	"github.com/golang_falcon_task/user-service/internal/logging"
	"github.com/golang_falcon_task/user-service/internal/metrics"
	"github.com/golang_falcon_task/user-service/internal/middleware"
	"github.com/golang_falcon_task/user-service/internal/service"
	"github.com/golang_falcon_task/user-service/internal/store"
	pb "github.com/golang_falcon_task/user-service/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
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
	metrics.StartMetricsServer(":9005")

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	userStore := store.NewPGUserStore(database)
	userService := service.NewUserService(userStore, log)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(log), // Logs all requests and responses
			middleware.MetricsInterceptor(),    // Captures Prometheus metrics
		),
	)

	pb.RegisterUserServiceServer(grpcServer, userService)

	// Enable server reflection
	reflection.Register(grpcServer)

	log.Println("UserService is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
