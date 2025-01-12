package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/manuelgcsousa/pokevote/internal/data"
)

type Server struct {
	Port     int
	Database *data.Database
	Tmpl     *template.Template
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: s.routes(),
	}

	// Server channel error
	serverErr := make(chan error, 1)

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Default().Info(fmt.Sprintf("Server listening on port %d...", s.Port))
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case <-quit:
		slog.Default().Info("Shutting down server...")
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			slog.Default().Error("Server error: " + err.Error())
		}
	}

	// Ensure resources tied to the context are released
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		slog.Default().Error("Server forced to shutdown: " + err.Error())
	}

	slog.Default().Info("Server exited gracefully.")
}

func (s *Server) routes() *chi.Mux {
	router := chi.NewRouter()

	// Minimal logger middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Default().Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
			next.ServeHTTP(w, r)
		})
	})

	// Setup routes
	router.Get("/", s.indexHandler)
	router.Get("/pokemon", s.pokemonHandler)
	router.Post("/vote", s.voteHandler)
	router.Get("/results", s.resultsHandler)

	return router
}
