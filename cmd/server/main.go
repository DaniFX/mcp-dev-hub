package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	ghclient "github.com/DaniFX/mcp-dev-hub/internal/github"
	"github.com/DaniFX/mcp-dev-hub/internal/mcp"
	"github.com/DaniFX/mcp-dev-hub/internal/mcp/tools"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// --- load GitHub accounts from env (single account quick setup) ---
	accounts := []ghclient.Account{}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		accounts = append(accounts, ghclient.Account{Alias: "default", Token: token})
	}
	// TODO: load additional accounts from configs/accounts.json

	mc := ghclient.NewMultiClient(accounts)

	// --- register MCP tools ---
	h := mcp.NewHandler()
	h.Register("list_repos", tools.ListReposHandler(mc))

	// --- HTTP routes ---
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/mcp", mcpHandler(h))

	log.Printf("mcp-dev-hub listening on :%s — tools: list_repos", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok","tools":["list_repos"]}`)
}

// mcpRequest represents a JSON-RPC 2.0 MCP request.
type mcpRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

// mcpResponse represents a JSON-RPC 2.0 MCP response.
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

		// Only handle tools/call for now
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
