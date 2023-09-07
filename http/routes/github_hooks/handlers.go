package github_hooks

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

func HandleCheckSuiteEvent(ctx context.Context, logger *slog.Logger, event *github.CheckSuiteEvent) error {

	logger = logger.With(
		slog.Group("event",
			slog.String("action", *event.Action),
			slog.String("repo", *event.Repo.FullName),
		),
	)

	logger.Info("Handle Check Suite Event")

	// //println(string(b))

	if *event.Action == "requested" || *event.Action == "rerequested" {

		// Action: requested
		// kick off a temporal workflow at this point to allow the creation of the suite.

		// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 385716, 41469183, "poc.pem")

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		// Use installation transport with client.
		client := github.NewClient(&http.Client{Transport: itr})

		var opts github.CreateCheckRunOptions
		opts.HeadSHA = *event.CheckSuite.HeadSHA
		opts.Name = "POC Run"
		opts.Status = github.String("queued")

		_, _, err = client.Checks.CreateCheckRun(ctx, *event.Repo.Owner.Login, *event.Repo.Name, opts)

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		logger.Info("Created Check Run!")

	}

	return nil

}

func HandleCheckRunEvent(ctx context.Context, logger *slog.Logger, event *github.CheckRunEvent) error {

	logger = logger.With(
		slog.Group("event",
			slog.String("action", *event.Action),
			slog.String("repo", *event.Repo.FullName),
		),
	)

	logger.Info("Handle Check Run Event")

	// //println(string(b))

	if *event.Action == "created" {

		// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 385716, 41469183, "poc.pem")

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		// Use installation transport with client.
		client := github.NewClient(&http.Client{Transport: itr})

		var output github.CheckRunOutput
		output.Title = github.String("Checks Completed")
		output.Summary = github.String("Organisaitonsal checks have been run.")

		var opts github.UpdateCheckRunOptions
		opts.Status = github.String("completed")
		opts.Conclusion = github.String("success")
		opts.Name = "POC Run"
		opts.Output = &output

		_, _, err = client.Checks.UpdateCheckRun(ctx, *event.Repo.Owner.Login, *event.Repo.Name, *event.CheckRun.ID, opts)

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		logger.Info("Completed Check Run!")

	}

	return nil

}
