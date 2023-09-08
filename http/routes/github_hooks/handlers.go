package github_hooks

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/go-github/github"
	"github.com/ross-mcdermott/github-app-temporal/workflows/activities"
	"github.com/ross-mcdermott/github-app-temporal/workflows/definitions"
	temporal "go.temporal.io/sdk/client"
)

func UnknownEvent(ctx context.Context, logger *slog.Logger, event interface{}) error {

	logger.Debug(fmt.Sprintf("Unknown event type '%T'. Ignoring.", event))

	return nil
}

func handlePullRequestEvent(ctx context.Context, logger *slog.Logger, event *github.PullRequestEvent) error {

	logger = logger.With(
		slog.Group("event",
			slog.String("action", *event.Action),
			slog.String("repo", *event.Repo.FullName),
		),
	)

	logger.Info("Handle PR Event")

	return nil
}

func handleCheckSuiteEvent(ctx context.Context, logger *slog.Logger, event *github.CheckSuiteEvent, client temporal.Client) error {

	logger = logger.With(
		slog.Group("event",
			slog.String("action", *event.Action),
			slog.String("repo", *event.Repo.FullName),
		),
	)

	logger.Info("Handle Check Suite Event")

	if *event.Action == "requested" || *event.Action == "rerequested" {

		workflowOptions := temporal.StartWorkflowOptions{
			TaskQueue: "default",
		}

		var params definitions.GithubCheckSuiteArgs
		params.Action = *event.Action
		params.HeadSHA = *event.CheckSuite.HeadSHA
		params.HeadBranch = *event.CheckSuite.HeadBranch
		params.Repo = activities.Repo{
			Name:       *event.Repo.Name,
			FullName:   *event.Repo.FullName,
			OwnerLogin: *event.Repo.Owner.Login,
		}

		workflowRun, err := client.ExecuteWorkflow(context.Background(), workflowOptions, definitions.GitHubCheckWorkflowDefinition, params)

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		logger.Debug(workflowRun.GetRunID())
	}

	return nil

}

func handleCheckRunEvent(ctx context.Context, logger *slog.Logger, event *github.CheckRunEvent, client temporal.Client) error {

	logger = logger.With(
		slog.Group("event",
			slog.String("action", *event.Action),
			slog.String("repo", *event.Repo.FullName),
		),
	)

	if event.CheckRun.ExternalID != nil {

		signalName := strings.ToLower("check_run:" + *event.Action)

		logger.Info("Signaling workflow", slog.String("signal", signalName))

		signal := definitions.CheckRunSignal{
			ID:         *event.CheckRun.ID,
			ExternalID: *event.CheckRun.ExternalID,
			Action:     *event.Action,
		}

		err := client.SignalWorkflow(context.Background(), *event.CheckRun.ExternalID, "", signalName, signal)
		if err != nil {
			return err
		}
	} else {
		logger.Warn("No external ID set. Ignoring.")
	}

	return nil
}
