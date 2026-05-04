package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
)

// IssueInfo is a simplified issue summary returned by MCP tools.
type IssueInfo struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	State     string   `json:"state"`
	Labels    []string `json:"labels"`
	Assignees []string `json:"assignees"`
	CreatedAt string   `json:"created_at"`
}

// CreateIssueOptions holds the parameters for creating a GitHub issue.
type CreateIssueOptions struct {
	Owner     string
	Repo      string
	Title     string
	Body      string
	Labels    []string
	Assignees []string
}

// CreateIssue creates a new issue in the given repo under the given account.
func (mc *MultiClient) CreateIssue(ctx context.Context, alias string, opts CreateIssueOptions) (*IssueInfo, error) {
	client, err := mc.For(alias)
	if err != nil {
		return nil, err
	}

	if opts.Title == "" {
		return nil, fmt.Errorf("issue title is required")
	}

	req := &github.IssueRequest{
		Title:     github.String(opts.Title),
		Body:      github.String(opts.Body),
		Labels:    &opts.Labels,
		Assignees: &opts.Assignees,
	}

	issue, _, err := client.Issues.Create(ctx, opts.Owner, opts.Repo, req)
	if err != nil {
		return nil, fmt.Errorf("error creating issue in %s/%s: %w", opts.Owner, opts.Repo, err)
	}

	labels := make([]string, 0, len(issue.Labels))
	for _, l := range issue.Labels {
		labels = append(labels, l.GetName())
	}
	assignees := make([]string, 0, len(issue.Assignees))
	for _, a := range issue.Assignees {
		assignees = append(assignees, a.GetLogin())
	}

	return &IssueInfo{
		Number:    issue.GetNumber(),
		Title:     issue.GetTitle(),
		URL:       issue.GetHTMLURL(),
		State:     issue.GetState(),
		Labels:    labels,
		Assignees: assignees,
		CreatedAt: issue.GetCreatedAt().String(),
	}, nil
}
