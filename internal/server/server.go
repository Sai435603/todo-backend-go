package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Sai435603/todo-backend-go/internal/app"
	"github.com/Sai435603/todo-backend-go/internal/handler"
	custommw "github.com/Sai435603/todo-backend-go/internal/middleware"
	routes "github.com/Sai435603/todo-backend-go/internal/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	app    *app.Application
	server *http.Server
}

func New(app *app.Application, h *handler.Handler, authHnd *handler.AuthHandler) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(custommw.RequestLogger(app.Logger))
	r.Use(custommw.CORS("*"))
	r.Use(custommw.SecurityHeaders)
	r.Use(custommw.ContentType)
	r.Use(custommw.RateLimiter(100, 200)) // 100 req/s per IP, burst of 200
	r.Use(custommw.Timeout(30 * time.Second))

	routes.Register(r, h, app.Config.JWTSecret)
	routes.RegisterAuthRoutes(r, authHnd)

	srv := &http.Server{
		Addr:         ":" + app.Config.App.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		app:    app,
		server: srv,
	}
}

func (s *Server) Start() error {

	s.app.Logger.Info(
		"HTTP server started",
		"address", s.server.Addr,
	)

	err := s.server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
