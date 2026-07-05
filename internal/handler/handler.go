package handler

import (
	"log/slog"

	"github.com/Sai435603/todo-backend-go/internal/service"
)

type Handler struct {
	logger       *slog.Logger
	service      *service.TodoService
	AuthHandler  *AuthHandler
}

func New(logger *slog.Logger, service *service.TodoService, authHandler *AuthHandler) *Handler {
	return &Handler{
		logger:      logger,
		service:     service,
		AuthHandler: authHandler,
	}
}
