package routes

import (
	"github.com/Sai435603/todo-backend-go/internal/handler"
	"github.com/go-chi/chi/v5"
)

func registerHealthRoutes(r chi.Router, h *handler.Handler) {
	r.Get("/health", h.Health)
}
