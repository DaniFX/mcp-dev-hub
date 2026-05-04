package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	ghclient "github.com/DaniFX/mcp-dev-hub/internal/github"
	"github.com/DaniFX/mcp-dev-hub/internal/mcp"
	"github.com/DaniFX/mcp-dev-hub/internal/mcp/tools"
)

// loadAccounts scans env vars with prefix GITHUB_TOKEN_<ALIAS>
// and builds the account list dynamically.
// Example:
//
//	GITHUB_TOKEN_PERSONAL=ghp_xxx  → alias "personal"
//	GITHUB_TOKEN_WORK=ghp_yyy      → alias "work"
func loadAccounts() []ghclient.Account {
	accounts := []ghclient.Account{}
	const prefix = "GITHUB_TOKEN_"
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 || parts[1] == "" {
			continue
		}
		alias := strings.ToLower(strings.TrimPrefix(parts[0], prefix))
		accounts = append(accounts, ghclient.Account{Alias: alias, Token: parts[1]})
		log.Printf("account loaded: %s", alias)
	}
	// fallback: plain GITHUB_TOKEN → alias "default"
	if len(accounts) == 0 {
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			accounts = append(accounts, ghclient.Account{Alias: "default", Token: token})
			log.Printf("account loaded: default (from GITHUB_TOKEN)")
		}
	}
	return accounts
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	accounts := loadAccounts()
	if len(accounts) == 0 {
		log.Println("[warn] no GitHub accounts configured — set GITHUB_TOKEN_<ALIAS> in .env")
	}

	mc := ghclient.NewMultiClient(accounts)

	// --- register MCP tools ---
	h := mcp.NewHandler()
	h.Register("list_repos", tools.ListReposHandler(mc))

	// --- HTTP routes ---
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/mcp", mcpHandler(h))

	log.Printf("mcp-dev-hub listening on :%s — accounts: %d — tools: list_repos", port, len(accounts))
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok","tools":["list_repos"]}`)
}

type mcpRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type mcpResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func mcpHandler(h *mcp.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req mcpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, req.ID, -32700, "parse error", err.Error())
			return
		}

		if req.Method != "tools/call" {
			writeError(w, req.ID, -32601, "method not found", req.Method)
			return
		}

		toolName, _ := req.Params["name"].(string)
		args, _ := req.Params["arguments"].(map[string]interface{})

		result, err := h.Call(toolName, args)
		if err != nil {
			writeError(w, req.ID, -32000, "tool error", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mcpResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		})
	}
}

func writeError(w http.ResponseWriter, id interface{}, code int, msg, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mcpResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: map[string]interface{}{
			"code":    code,
			"message": msg,
			"data":    data,
		},
	})
}
