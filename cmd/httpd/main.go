package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/ross-mcdermott/github-app-temporal/http/handlers"
	"github.com/ross-mcdermott/github-app-temporal/http/webhooks"
	"go.temporal.io/sdk/client"
)

func main() {

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	//logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// set default
	slog.SetDefault(logger)

	// Create a Temporal Client to communicate with the Temporal Cluster.
	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client", err)
	}
	defer temporalClient.Close()

	r := chi.NewRouter()

	r.Get("/security.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO"))
	})

	// Register health check endpoint
	const health_route = "/healthz"
	healthCheck := handlers.NewHealthHandler(logger)
	healthCheck.Register(r, health_route)

	// Register the github hooks
	const github_route = "/hooks/github"
	ghHooks := webhooks.NewGithubHandler(logger, temporalClient, "0695679902")
	ghHooks.Register(r, github_route)

	// Now start the API
	http.ListenAndServe(":3000", r)
}
