package handler

import (
	"log/slog"

	"github.com/Sai435603/todo-backend-go/internal/service"
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
