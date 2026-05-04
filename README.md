# mcp-dev-hub

> Remote MCP Server for multi-account GitHub management, CI/CD automation, and automatic documentation.

## Architecture

```
mcp-dev-hub/
├── cmd/server/          # Entry point
├── internal/
│   ├── mcp/             # Core MCP handler (tools, resources)
│   ├── github/          # Multi-account GitHub client
│   ├── firestore/       # Firestore state & metadata layer
│   ├── auth/            # API key middleware
│   ├── docs/            # Automatic documentation generator
│   └── cicd/            # CI/CD automation layer
├── configs/             # Account configs (gitignored)
├── Dockerfile
└── .gitignore
```

## Stack

| Component | Technology |
|---|---|
| Language | Go |
| Runtime | Google Cloud Run |
| Database | Firestore |
| Storage | Firebase Storage |
| Auth (v1) | API Key via Secret Manager |
| GitHub | Multi-account via token routing |

## MCP Tools (planned)

### GitHub
- `list_repos` — list repos across all configured accounts
- `get_repo` — get repo details
- `create_pr` — open a PR on a target account
- `get_pr_status` — review PR status and checks

### Docs
- `generate_docs` — generate documentation from code/README
- `get_doc_summary` — summarize a repo's documentation
- `update_readme` — auto-update README on push

### CI/CD
- `get_pipeline_status` — get GitHub Actions build status
- `analyze_failure` — diagnose failed builds
- `trigger_workflow` — trigger a workflow dispatch
- `get_deploy_history` — fetch deploy history for a project

## Getting Started

```bash
# Clone the repo
git clone https://github.com/DaniFX/mcp-dev-hub.git
cd mcp-dev-hub

# Install dependencies
go mod tidy

# Run locally
go run ./cmd/server
```

## Deployment

```bash
# Build and push Docker image
gcloud builds submit --tag gcr.io/YOUR_PROJECT/mcp-dev-hub

# Deploy to Cloud Run
gcloud run deploy mcp-dev-hub \
  --image gcr.io/YOUR_PROJECT/mcp-dev-hub \
  --platform managed \
  --region europe-west1 \
  --allow-unauthenticated
```

## License

MIT
