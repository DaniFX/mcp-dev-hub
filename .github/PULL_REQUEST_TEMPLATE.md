## 📋 Descrizione

<!-- Descrivi brevemente cosa fa questa PR e perché è necessaria -->

## 🔗 Issue collegata

Closes #<!-- numero issue -->

## 🧩 Tipo di modifica

- [ ] 🐛 Bug fix
- [ ] ✨ Nuova feature
- [ ] ♻️ Refactor
- [ ] 📝 Documentazione
- [ ] ⚙️ CI/CD / Config
- [ ] 🔒 Security

## 🛠️ Layer coinvolti

- [ ] `mcp/` — Core handler
- [ ] `github/` — Multi-account client
- [ ] `auth/` — Middleware
- [ ] `firestore/` — Database layer
- [ ] `docs/` — Documentazione automatica
- [ ] `cicd/` — CI/CD automation
- [ ] `cmd/server` — Entry point
- [ ] `Dockerfile` / Deploy
- [ ] `.devcontainer` / Config

## ✅ Checklist

- [ ] Il codice compila senza errori (`go build ./...`)
- [ ] I test passano (`go test ./...`)
- [ ] Il linter non ha errori (`golangci-lint run ./...`)
- [ ] Ho aggiornato il README se necessario
- [ ] Ho aggiornato/aggiunto i tool MCP nel manifest se necessario
- [ ] Non ho committato secrets o token

## 🧪 Come testare

<!-- Descrivi i passi per verificare che la modifica funzioni -->

1. Apri il Codespace su questo branch
2. Avvia il server (`go run ./cmd/server`)
3. Usa `api/test.http` per testare l'endpoint
4. ...

## 📸 Screenshot / Output (se applicabile)

<!-- Incolla qui output del terminale o screenshot se utile -->
