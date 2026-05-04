package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/mcp", mcpHandler)

	log.Printf("mcp-dev-hub server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "ok"}`)
}

func mcpHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: route MCP JSON-RPC 2.0 requests
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"jsonrpc": "2.0", "result": "mcp-dev-hub ready", "id": null}`)
}
