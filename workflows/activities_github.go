package workflows

import (
	"context"
	"log/slog"

	"github.com/google/go-github/github"
)

// YourActivityObject is the struct that maintains shared state across Activities.
// If the Worker crashes this Activity object loses its state.
type GitHubActivities struct {
	Client *github.Client
	Logger *slog.Logger
}

type Repo struct {
	Name       string
	FullName   string
	OwnerLogin string
}

type CreateCheckRunActivityArgs struct {
	Name       string
	HeadSha    string
	ExternalId string
	Repo       Repo
}

type CheckRunStatus struct {
	ID         int64
	ExternalID string
	Status     string
}

func (a *GitHubActivities) CreateCheckRun(ctx context.Context, param CreateCheckRunActivityArgs) (*CheckRunStatus, error) {

	var opts github.CreateCheckRunOptions
	opts.ExternalID = &param.ExternalId
	opts.HeadSHA = param.HeadSha
	opts.Name = param.Name
	opts.Status = github.String("queued")

	checkRun, _, err := a.Client.Checks.CreateCheckRun(ctx, param.Repo.OwnerLogin, param.Repo.Name, opts)

	if err != nil {
		a.Logger.Error(err.Error())
		return nil, err
	}

	a.Logger.Info("Created Check Run!")

	result := &CheckRunStatus{
		ID:         checkRun.GetID(),
		ExternalID: *checkRun.ExternalID,
		Status:     *checkRun.Status,
	}

	// Return the results back to the Workflow Execution.
	return result, nil
}

type UpdateCheckRunActivityArgs struct {
	ID         int64
	Name       string
	Repo       Repo
	Status     string
	Conculsion string
}

func (a *GitHubActivities) UpdateCheckRun(ctx context.Context, param UpdateCheckRunActivityArgs) (*CheckRunStatus, error) {

	// var output github.CheckRunOutput
	// output.Title = github.String("Checks Completed")
	// output.Summary = github.String("Organisaitonsal checks have been run.")

	var opts github.UpdateCheckRunOptions
	opts.Name = param.Name
	opts.Status = &param.Status // github.String("completed")
	if param.Conculsion != "" {
		opts.Conclusion = &param.Conculsion // github.String("success")
	}

	//opts.Output = &output

	checkRun, _, err := a.Client.Checks.UpdateCheckRun(ctx, param.Repo.OwnerLogin, param.Repo.Name, param.ID, opts)

	if err != nil {
		a.Logger.Error(err.Error())
		return nil, err
	}

	a.Logger.Info("Updated Check Run")

	result := &CheckRunStatus{
		ID:         checkRun.GetID(),
		ExternalID: *checkRun.ExternalID,
		Status:     *checkRun.Status,
	}

	// Return the results back to the Workflow Execution.
	return result, nil

}
