package tools

import (
	"context"
	"encoding/json"
	"fmt"

	ghclient "github.com/DaniFX/mcp-dev-hub/internal/github"
)

// ListReposHandler returns an MCP ToolFunc for the list_repos tool.
//
// Expected params:
//
//	"account" (string, required) — alias of the GitHub account
//	"type"    (string, optional) — "all"|"public"|"private"|"owner" (default: "owner")
//	"page"    (float64, optional) — page number for pagination (default: 1)
func ListReposHandler(mc *ghclient.MultiClient) func(map[string]interface{}) (interface{}, error) {
	return func(params map[string]interface{}) (interface{}, error) {
		// --- parse params ---
		account, ok := params["account"].(string)
		if !ok || account == "" {
			return nil, fmt.Errorf("missing required param: 'account'. Available accounts: %v", mc.Aliases())
		}

		repoType, _ := params["type"].(string)

		page := 1
		if p, ok := params["page"].(float64); ok && p > 0 {
			page = int(p)
		}

		// --- call GitHub ---
		repos, err := mc.ListRepos(context.Background(), account, repoType, page)
		if err != nil {
			return nil, err
		}

		// --- build MCP response ---
		data, err := json.Marshal(map[string]interface{}{
			"account": account,
			"type":    repoType,
			"page":    page,
			"count":   len(repos),
			"repos":   repos,
		})
		if err != nil {
			return nil, fmt.Errorf("error serializing response: %w", err)
		}

		return json.RawMessage(data), nil
	}
}
