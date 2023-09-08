package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type GithubCheckSuiteArgs struct {
	Action     string
	HeadSHA    string
	HeadBranch string
	Repo       Repo
}

type Repo struct {
	Name       string
	FullName   string
	OwnerLogin string
}

type GitHubCheckWorkflowDefinitionResult struct {
	ExternalID string
	Status     string
}

type CheckRunSignal struct {
	ID         int64
	ExternalID string
	Action     string
}

func GitHubCheckWorkflowDefinition(ctx workflow.Context, param GithubCheckSuiteArgs) (*GitHubCheckWorkflowDefinitionResult, error) {

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	createParams := CreateCheckRunActivityArgs{
		Name:       "Sample Check",
		HeadSha:    param.HeadSHA,
		Repo:       param.Repo,
		ExternalId: workflow.GetInfo(ctx).WorkflowExecution.ID,
	}

	// Use a nil struct pointer to call Activities that are part of a struct.
	var a *GitHubActivities
	// Execute the Activity and wait for the result.
	var createParamsResult *CheckRunStatus
	err := workflow.ExecuteActivity(ctx, a.CreateCheckRun, createParams).Get(ctx, &createParamsResult)
	if err != nil {
		return nil, err
	}

	var checkRunCreated CheckRunSignal
	signalChan := workflow.GetSignalChannel(ctx, "check_run:created")
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(signalChan, func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &checkRunCreated)
	})
	selector.Select(ctx)

	// we received it now.
	updateParams := UpdateCheckRunActivityArgs{
		ID:     createParamsResult.ID,
		Name:   createParams.Name,
		Repo:   param.Repo,
		Status: "in_progress",
		// no conculsuion at this point
	}

	var updateParamsResult *CheckRunStatus
	err = workflow.ExecuteActivity(ctx, a.UpdateCheckRun, updateParams).Get(ctx, &updateParamsResult)
	if err != nil {
		return nil, err
	}

	// sleep 10 seconds for fun
	workflow.Sleep(ctx, 10*time.Second)

	// we received it now.
	updateParams = UpdateCheckRunActivityArgs{
		ID:         createParamsResult.ID,
		Name:       createParams.Name,
		Repo:       param.Repo,
		Status:     "completed",
		Conculsion: "success",
	}

	err = workflow.ExecuteActivity(ctx, a.UpdateCheckRun, updateParams).Get(ctx, &updateParamsResult)
	if err != nil {
		return nil, err
	}

	// Make the results of the Workflow Execution available.
	workflowResult := &GitHubCheckWorkflowDefinitionResult{
		ExternalID: updateParamsResult.ExternalID,
		Status:     updateParamsResult.Status,
	}

	return workflowResult, nil
}
