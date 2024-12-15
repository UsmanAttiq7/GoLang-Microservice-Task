package db

import (
	"context"
	"fmt"

	"github.com/golang_falcon_task/ride-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB initializes a database connection pool.
func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return pool, nil
}
