package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/Sai435603/todo-backend-go/internal/app"
	"github.com/Sai435603/todo-backend-go/internal/handler"
	routes "github.com/Sai435603/todo-backend-go/internal/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	app    *app.Application
	server *http.Server
}

func New(app *app.Application, h *handler.Handler) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	routes.Register(r, h)

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
