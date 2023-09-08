package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	ghhooks "github.com/ross-mcdermott/github-app-temporal/http/routes/github_hooks"
	"github.com/ross-mcdermott/github-app-temporal/http/routes/health"
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
	// A Temporal Client is a heavyweight object that should be created just once per process.
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

	// Register the routes to serve
	health.RegisterRoutes(r, logger)
	ghhooks.RegisterRoutes(r, logger, temporalClient)

	// Now start the API
	http.ListenAndServe(":3000", r)
}
