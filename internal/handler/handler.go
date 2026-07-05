package handler

import (
	"log/slog"

	"github.com/Sai435603/todo-backend-go/internal/config"
	"github.com/Sai435603/todo-backend-go/internal/database/sqlc"
	"github.com/Sai435603/todo-backend-go/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	logger  *slog.Logger
	service *service.TodoService
}

func New(logger *slog.Logger, service *service.TodoService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func NewAuthHandler(config *config.Config, db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		AuthConfig: config.GoogleOauthConfig,
		Cookie:     config.Cookie,
		JWTSecret:  config.JWTSecret,
		Queries:    sqlc.New(db),
	}
}
