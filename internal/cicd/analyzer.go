package cicd

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v62/github"
)

// Analyzer inspects failed workflow runs and extracts failure reasons.
type Analyzer struct {
	client *github.Client
}

// NewAnalyzer creates a new CI failure analyzer.
func NewAnalyzer(client *github.Client) *Analyzer {
	return &Analyzer{client: client}
}

// AnalyzeFailure returns a human-readable summary of a failed run.
func (a *Analyzer) AnalyzeFailure(ctx context.Context, owner, repo string, runID int64) (string, error) {
	jobs, _, err := a.client.Actions.ListWorkflowJobsAttempt(ctx, owner, repo, runID, 0, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list jobs: %w", err)
	}

	var sb strings.Builder
	for _, job := range jobs.Jobs {
		if job.GetConclusion() == "failure" {
			sb.WriteString(fmt.Sprintf("❌ Job '%s' failed\n", job.GetName()))
			for _, step := range job.Steps {
				if step.GetConclusion() == "failure" {
					sb.WriteString(fmt.Sprintf("  └─ Step: %s\n", step.GetName()))
				}
			}
		}
	}

	if sb.Len() == 0 {
		return "No failures detected in this run.", nil
	}
	return sb.String(), nil
}
