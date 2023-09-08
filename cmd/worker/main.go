package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/ross-mcdermott/github-app-temporal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

	internalWorker := worker.New(temporalClient, "default", worker.Options{})

	internalWorker.RegisterWorkflow(workflows.GitHubCheckWorkflowDefinition)

	// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 385716, 41469183, "poc.pem")

	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr})

	// Struct contains all the activities related to github
	ghactivities := &workflows.GitHubActivities{
		Client: client,
		Logger: logger,
	}

	// Register the activities that are available
	internalWorker.RegisterActivity(ghactivities)

	// Run the Worker
	err = internalWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}

}
