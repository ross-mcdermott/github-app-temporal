package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	ghhooks "github.com/ross-mcdermott/github-app-temporal/http/routes/github_hooks"
	"github.com/ross-mcdermott/github-app-temporal/http/routes/health"
)

func main() {
	r := chi.NewRouter()

	r.Get("/security.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO"))
	})

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	//logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// Register the routes to serve
	health.RegisterRoutes(r, logger)
	ghhooks.RegisterRoutes(r, logger)

	http.ListenAndServe(":3000", r)
}
