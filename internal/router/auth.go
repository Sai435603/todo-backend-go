package routes

import (
	"github.com/Sai435603/todo-backend-go/internal/handler"
	"github.com/go-chi/chi/v5"
)

func registerAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Get("/auth/google", h.HandleOAuthLogin)
	r.Get("/auth/google/callback", h.OAuthCallbackHandler)
}
