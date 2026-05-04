package cicd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
)

// PipelineService handles CI/CD pipeline operations via GitHub Actions.
type PipelineService struct {
	client *github.Client
}

// NewPipelineService creates a new PipelineService for a given GitHub client.
func NewPipelineService(client *github.Client) *PipelineService {
	return &PipelineService{client: client}
}

// GetStatus returns the latest workflow run status for a repo.
func (p *PipelineService) GetStatus(ctx context.Context, owner, repo string) (string, error) {
	runs, _, err := p.client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow runs: %w", err)
	}
	if runs.GetTotalCount() == 0 {
		return "no runs found", nil
	}
	latest := runs.WorkflowRuns[0]
	return fmt.Sprintf("status=%s conclusion=%s", latest.GetStatus(), latest.GetConclusion()), nil
}

// TriggerWorkflow dispatches a workflow_dispatch event.
func (p *PipelineService) TriggerWorkflow(ctx context.Context, owner, repo, workflow, ref string) error {
	event := github.CreateWorkflowDispatchEventRequest{Ref: ref}
	_, err := p.client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, workflow, event)
	return err
}
