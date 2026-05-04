package auth

import (
	"net/http"
	"os"
)

// APIKeyMiddleware validates requests using a static API key from env.
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := os.Getenv("MCP_API_KEY")
		provided := r.Header.Get("X-API-Key")
		if expected != "" && provided != expected {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
