#!/bin/bash
set -e

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "==> Go version: $GO_VERSION"

echo "==> [1/5] Installing Go tools..."
# Pinned versions compatible with Go 1.24.x
go install golang.org/x/tools/gopls@v0.19.0
go install github.com/go-delve/delve/cmd/dlv@v1.23.1
go install golang.org/x/tools/cmd/goimports@v0.29.0
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

echo "==> [2/5] Installing GitHub CLI..."
if ! command -v gh &> /dev/null; then
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
  sudo apt-get update -qq && sudo apt-get install -y gh
else
  echo "    gh already installed: $(gh --version | head -1)"
fi

echo "==> [3/5] Installing Google Cloud CLI..."
if ! command -v gcloud &> /dev/null; then
  curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg 2>/dev/null
  echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list > /dev/null
  sudo apt-get update -qq && sudo apt-get install -y google-cloud-cli
else
  echo "    gcloud already installed: $(gcloud --version | head -1)"
fi

echo "==> [4/5] Downloading Go dependencies..."
go mod tidy

echo "==> [5/5] Copying example configs..."
if [ ! -f configs/accounts.json ]; then
  cp configs/accounts.example.json configs/accounts.json
  echo "    configs/accounts.json created (fill in your tokens)"
fi

if [ ! -f .env ]; then
  cp .env.example .env
  echo "    .env created (fill in your values)"
fi

echo ""
echo "✅ mcp-dev-hub devcontainer ready!"
echo "   Run:   go run ./cmd/server"
echo "   Debug: F5 in VSCode"
echo "   Test:  go test ./..."
echo "   Lint:  golangci-lint run ./..."
