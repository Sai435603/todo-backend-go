package app

import (
	"log/slog"

	"github.com/Sai435603/todo-backend-go/internal/config"
	"github.com/Sai435603/todo-backend-go/internal/database"
	"github.com/Sai435603/todo-backend-go/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *pgxpool.Pool
}

func New() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg)
	log.Info("Application starting...")
	log.Info(
		"Configuration loaded",
		"port", cfg.App.Port,
		"environment", cfg.App.Env,
	)
	db, err := database.New(cfg, log)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return nil, err
	}
	// defer db.Close()
	log.Info("database connection established")

	return &Application{
		Config: cfg,
		Logger: log,
		DB:     db,
	}, nil
}
