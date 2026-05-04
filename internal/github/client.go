package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

// Account represents a GitHub account with its access token.
type Account struct {
	Alias string `json:"alias"`
	Token string `json:"token"`
}

// RepoInfo is a simplified repo summary returned by MCP tools.
type RepoInfo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	URL         string `json:"url"`
	Language    string `json:"language"`
	Stars       int    `json:"stars"`
	UpdatedAt   string `json:"updated_at"`
}

// MultiClient manages multiple GitHub account clients.
type MultiClient struct {
	clients map[string]*github.Client
}

// NewMultiClient creates clients for each configured account.
func NewMultiClient(accounts []Account) *MultiClient {
	mc := &MultiClient{clients: make(map[string]*github.Client)}
	for _, acc := range accounts {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: acc.Token})
		tc := oauth2.NewClient(context.Background(), ts)
		mc.clients[acc.Alias] = github.NewClient(tc)
	}
	return mc
}

// For returns the GitHub client for a given account alias.
func (mc *MultiClient) For(alias string) (*github.Client, error) {
	c, ok := mc.clients[alias]
	if !ok {
		return nil, fmt.Errorf("account '%s' not configured", alias)
	}
	return c, nil
}

// Aliases returns the list of configured account aliases.
func (mc *MultiClient) Aliases() []string {
	keys := make([]string, 0, len(mc.clients))
	for k := range mc.clients {
		keys = append(keys, k)
	}
	return keys
}

// ListRepos returns repos for the given account, with optional type filter.
// repoType: "all" | "public" | "private" | "owner" (default: "owner")
func (mc *MultiClient) ListRepos(ctx context.Context, alias string, repoType string, page int) ([]RepoInfo, error) {
	client, err := mc.For(alias)
	if err != nil {
		return nil, err
	}

	if repoType == "" {
		repoType = "owner"
	}
	if page < 1 {
		page = 1
	}

	opts := &github.RepositoryListByAuthenticatedUserOptions{
		Type: repoType,
		Sort: "updated",
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: 30,
		},
	}

	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error listing repos for account '%s': %w", alias, err)
	}

	result := make([]RepoInfo, 0, len(repos))
	for _, r := range repos {
		result = append(result, RepoInfo{
			Name:        r.GetName(),
			FullName:    r.GetFullName(),
			Description: r.GetDescription(),
			Private:     r.GetPrivate(),
			URL:         r.GetHTMLURL(),
			Language:    r.GetLanguage(),
			Stars:       r.GetStargazersCount(),
			UpdatedAt:   r.GetUpdatedAt().String(),
		})
	}
	return result, nil
}
