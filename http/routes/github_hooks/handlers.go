package github_hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v55/github"
)

func HandleCheckSuiteEvent(ctx context.Context, logger *slog.Logger, b []byte, deliveryID string) error {

	logger.Info("Handle Check Suite")

	//println(string(b))

	var c CheckSuiteWebhook
	decoder := json.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&c)

	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("Action '%s' Repository '%s'", c.Action, c.Repository.FullName))

	if c.Action == "requested" || c.Action == "rerequested" {
		// Action: requested
		// kick off a temporal workflow at this point to allow the creation of the suite.

		// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 385716, 41469183, "poc.pem")

		// Or for endpoints that require JWT authentication
		// itr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, 1, "2016-10-19.private-key.pem")

		if err != nil {
			// Handle error.
		}

		// Use installation transport with client.
		client := github.NewClient(&http.Client{Transport: itr})

		var opts github.CreateCheckRunOptions
		opts.HeadSHA = c.CheckSuite.HeadSha
		opts.Name = "POC Run"

		_, _, err = client.Checks.CreateCheckRun(ctx, c.Repository.Owner.Login, c.Repository.Name, opts)

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		logger.Info("Created Check Run!")

	}

	return nil

}

func HandleCheckRunEvent(ctx context.Context, logger *slog.Logger, b []byte, deliveryID string) error {

	logger.Info("Handle Check Run")

	//println(string(b))

	var c CheckRunWebhook
	decoder := json.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&c)

	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("Action '%s' Repository '%s'", c.Action, c.Repository.FullName))

	if c.Action == "created" {

		// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 385716, 41469183, "poc.pem")

		// Or for endpoints that require JWT authentication
		// itr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, 1, "2016-10-19.private-key.pem")

		if err != nil {
			// Handle error.
		}

		// Use installation transport with client.
		client := github.NewClient(&http.Client{Transport: itr})

		var opts github.UpdateCheckRunOptions
		opts.Status = github.String("completed")
		opts.Conclusion = github.String("success")
		opts.Name = "POC Run"

		_, _, err = client.Checks.UpdateCheckRun(ctx, c.Repository.Owner.Login, c.Repository.Name, c.CheckRun.ID, opts)

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		logger.Info("Completed Check Run!")

	}

	return nil

}
