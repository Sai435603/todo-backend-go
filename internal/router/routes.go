package routes

import (
	"github.com/Sai435603/todo-backend-go/internal/handler"
	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router, h *handler.Handler) {
	r.Route("/api/v1", func(r chi.Router) {
		registerHealthRoutes(r, h)
		registerTodoRoutes(r, h)
	})
}
