#!/bin/bash
set -e

echo "==> Installing Go tools..."
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

echo "==> Downloading Go dependencies..."
go mod tidy

echo "==> Copying example configs..."
if [ ! -f configs/accounts.json ]; then
  cp configs/accounts.example.json configs/accounts.json
  echo "    configs/accounts.json created from example (fill in your tokens)"
fi

if [ ! -f .env ]; then
  cp .env.example .env
  echo "    .env created from example (fill in your values)"
fi

echo ""
echo "✅ mcp-dev-hub devcontainer ready!"
echo "   Run server:  go run ./cmd/server"
echo "   Debug:       F5 in VSCode or 'dlv debug ./cmd/server'"
echo "   Test:        go test ./..."
echo "   Lint:        golangci-lint run ./..."
