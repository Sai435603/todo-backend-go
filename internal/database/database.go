package database

import (
	"context"
	"log/slog"

	"github.com/Sai435603/todo-backend-go/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	ctx := context.Background()

	logger.Info("Connecting to PostgreSQL...")

	pool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		logger.Error("Failed to create connection pool", "error", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Error("Failed to ping PostgreSQL", "error", err)
		return nil, err
	}

	logger.Info("Connected to PostgreSQL")

	return pool, nil
}
