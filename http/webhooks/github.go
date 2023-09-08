package webhooks

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-github/github"
	"github.com/ross-mcdermott/github-app-temporal/workflows"
	temporal "go.temporal.io/sdk/client"
)

func NewGithubHandler(logger *slog.Logger, client temporal.Client, webhookSecret string) *GithubHooks {

	// Create new instance of the hooks handler
	handler := &GithubHooks{
		logger:        logger,
		client:        client,
		webhookSecret: webhookSecret,
	}

	return handler
}

type GithubHooks struct {
	logger        *slog.Logger
	client        temporal.Client
	webhookSecret string
}

func (a *GithubHooks) Register(router chi.Router, route string) {

	a.logger.Debug(fmt.Sprintf("Registering '%s' route.", route))

	router.Group(func(r chi.Router) {
		r.Post(route, a.process)
	})

}

func (a *GithubHooks) process(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	r = r.WithContext(ctx)

	payload, err := github.ValidatePayload(r, []byte(a.webhookSecret))
	if err != nil {
		a.logger.Error("Unable to validate payload")
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		a.logger.Error("Unable to parse payload")
		return
	}

	switch event := event.(type) {
	case *github.CheckSuiteEvent:
		err = a.webhook_check_suite(ctx, event)
	case *github.CheckRunEvent:
		err = a.webhook_check_run(ctx, event)
	case *github.PullRequestEvent:
		err = a.webhook_pull_request(ctx, event)
	default:
		a.logger.Debug(fmt.Sprintf("Unknown event type '%T'. Ignoring.", event))
	}

}

func (a *GithubHooks) webhook_pull_request(ctx context.Context, event *github.PullRequestEvent) error {

	a.logger.Info("Github 'Pull Request' Webhook Received",
		slog.String("action", *event.Action),
		slog.String("repo", *event.Repo.FullName),
	)

	return nil
}

func (a *GithubHooks) webhook_check_suite(ctx context.Context, event *github.CheckSuiteEvent) error {

	a.logger.Info("Github 'Check Suite' Webhook Received",
		slog.String("action", *event.Action),
		slog.String("repo", *event.Repo.FullName),
	)

	if *event.Action == "requested" || *event.Action == "rerequested" {

		workflowOptions := temporal.StartWorkflowOptions{
			TaskQueue: "default",
		}

		var params workflows.GithubCheckSuiteArgs
		params.Action = *event.Action
		params.HeadSHA = *event.CheckSuite.HeadSHA
		params.HeadBranch = *event.CheckSuite.HeadBranch
		params.Repo = workflows.Repo{
			Name:       *event.Repo.Name,
			FullName:   *event.Repo.FullName,
			OwnerLogin: *event.Repo.Owner.Login,
		}

		workflowRun, err := a.client.ExecuteWorkflow(context.Background(), workflowOptions, workflows.GitHubCheckWorkflowDefinition, params)

		if err != nil {
			return err
		}

		a.logger.Info("Started workflow execution",
			slog.String("workflow_id", workflowRun.GetID()),
			slog.String("workflow_run_id", workflowRun.GetRunID()),
		)
	}

	return nil

}

func (a *GithubHooks) webhook_check_run(ctx context.Context, event *github.CheckRunEvent) error {

	if event.CheckRun.ExternalID != nil {

		signalName := strings.ToLower("check_run:" + *event.Action)

		a.logger.Info("Signaling workflow",
			slog.String("workflow_id", *event.CheckRun.ExternalID),
			slog.String("action", *event.Action),
			slog.String("signal", signalName),
		)

		signal := workflows.CheckRunSignal{
			ID:         *event.CheckRun.ID,
			ExternalID: *event.CheckRun.ExternalID,
			Action:     *event.Action,
		}

		err := a.client.SignalWorkflow(context.Background(), *event.CheckRun.ExternalID, "", signalName, signal)

		if err != nil {
			return err
		}

	} else {
		a.logger.Debug("No external ID set. Ignoring.")
	}

	return nil
}
