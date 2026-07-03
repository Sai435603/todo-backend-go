package routes

import (
	"github.com/Sai435603/todo-backend-go/internal/handler"
	"github.com/go-chi/chi/v5"
)

func registerTodoRoutes(r chi.Router, h *handler.Handler) {
	r.Route("/todos", func(r chi.Router) {

		r.Get("/", h.GetTodos)
		r.Post("/", h.CreateTodo)

		r.Get("/{id}", h.GetTodo)
		r.Put("/{id}", h.UpdateTodo)
		r.Delete("/{id}", h.DeleteTodo)

		r.Get("/completed", h.GetCompletedTodos)
		r.Get("/pending", h.GetPendingTodos)

		r.Patch("/{id}/complete", h.MarkTodoCompleted)
		r.Patch("/{id}/pending", h.MarkTodoPending)

		r.Get("/search", h.SearchTodos)
		r.Get("/range", h.GetTodosByDateRange)
	})
}
