package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	config.MinConns = 2
	config.MaxConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 10 * time.Minute
	config.HealthCheckPeriod = time.Minute

	config.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL with connection pool!")

	return pool, err
}
