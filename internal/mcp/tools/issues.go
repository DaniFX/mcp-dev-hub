package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ghclient "github.com/DaniFX/mcp-dev-hub/internal/github"
)

// CreateIssueHandler returns an MCP ToolFunc for the create_issue tool.
//
// Expected params:
//
//	"account"   (string, required) — alias of the GitHub account
//	"repo"      (string, required) — "owner/repo" format
//	"title"     (string, required) — issue title
//	"body"      (string, optional) — issue body in Markdown
//	"labels"    ([]string, optional) — label names (created if they don't exist)
//	"assignees" ([]string, optional) — GitHub usernames to assign
func CreateIssueHandler(mc *ghclient.MultiClient) func(map[string]interface{}) (interface{}, error) {
	return func(params map[string]interface{}) (interface{}, error) {
		// --- parse required params ---
		account, ok := params["account"].(string)
		if !ok || account == "" {
			return nil, fmt.Errorf("missing required param: 'account'. Available: %v", mc.Aliases())
		}

		repoParam, ok := params["repo"].(string)
		if !ok || repoParam == "" {
			return nil, fmt.Errorf("missing required param: 'repo' (format: 'owner/repo')")
		}
		parts := strings.SplitN(repoParam, "/", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid 'repo' format: expected 'owner/repo', got '%s'", repoParam)
		}
		owner, repo := parts[0], parts[1]

		title, ok := params["title"].(string)
		if !ok || title == "" {
			return nil, fmt.Errorf("missing required param: 'title'")
		}

		// --- parse optional params ---
		body, _ := params["body"].(string)

		labels := toStringSlice(params["labels"])
		assignees := toStringSlice(params["assignees"])

		// --- call GitHub ---
		issue, err := mc.CreateIssue(context.Background(), account, ghclient.CreateIssueOptions{
			Owner:     owner,
			Repo:      repo,
			Title:     title,
			Body:      body,
			Labels:    labels,
			Assignees: assignees,
		})
		if err != nil {
			return nil, err
		}

		// --- build MCP response ---
		data, err := json.Marshal(map[string]interface{}{
			"created": true,
			"issue":   issue,
		})
		if err != nil {
			return nil, fmt.Errorf("error serializing response: %w", err)
		}

		return json.RawMessage(data), nil
	}
}

// toStringSlice safely converts interface{} to []string.
// Accepts []interface{} (from JSON unmarshalling) or []string.
func toStringSlice(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	switch val := v.(type) {
	case []string:
		return val
	case []interface{}:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return []string{}
}
