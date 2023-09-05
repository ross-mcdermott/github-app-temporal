package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ross-mcdermott/github-app-temporal/internal/routes/github_hooks"
	"github.com/ross-mcdermott/github-app-temporal/internal/routes/health"
)

func main() {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)

	r.Get("/security.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO"))
	})

	// Register routes as groups to allow differing middleware to be
	// associated with each set.
	r.Group(health.Routes)       // Health probes
	r.Group(github_hooks.Routes) // Github Hooks

	http.ListenAndServe(":3000", r)
}
