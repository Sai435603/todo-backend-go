package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sai435603/todo-backend-go/internal/handler"
	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router, h *handler.Handler) {
	r.Route("/api/v1", func(r chi.Router) {
		registerHealthRoutes(r, h)
		registerTodoRoutes(r, h)
	})

	// Serve static files from the "static" directory
	staticDir := "static"
	fileServer := http.FileServer(http.Dir(staticDir))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Check if file exists, if not serve index.html
		fullPath := filepath.Join(staticDir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		fileServer.ServeHTTP(w, r)
	})
}
