# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

projman is a Go CLI tool for managing development projects. It integrates with tmux and VS Code to launch sessions, supports local project discovery, GitHub remote repos, project templating, and cloning.

## Build & Development Commands

Task runner is `mise` (configured in `mise.toml`).

| Task | Command | Purpose |
|------|---------|---------|
| `mise run dev` | `go run .` | Run in development |
| `mise run build` | `go build .` | Build binary |
| `mise run build:install` | Build + install to `~/.local/bin/` | Build and install |
| `mise run test` | `bash scripts/test.sh` | Run all tests (uses gotestsum if available, else go test) |
| `mise run lint:fmt` | `gofmt -l .` | Check formatting (CI fails on unformatted code) |
| `mise run lint:vet` | `go vet ./...` | Run Go vet |
| `mise run lint:static` | `golangci-lint run` | Run aggregated linters |
| `mise run lint:deps` | `go mod tidy && git diff --exit-code go.mod go.sum` | Check dependency tidiness |
| `mise run clean` | `bash scripts/clean.sh` | Remove built binary |

Direct Go commands also work:
- Test all: `go test ./...`
- Test single package: `go test ./internal/services`
- Test single test: `go test ./internal/services -run TestFunctionName`
- Format fix: `gofmt -w .`

**Important** Before committing any work in this project, run `clint` to check the CI

## Architecture

Three-layer architecture: **Controllers** (CLI handlers) -> **Services** (business logic) -> **Repositories** (config/data).

```
internal/
  controllers/   # One file per CLI command (open, remote, new, clone, rm, etc.)
  services/      # Business logic (projects, github, tmux, vscode, select, health, etc.)
  repositories/  # Config file management
  ui/            # Spinner component
```

**Dependency flow**: `main.go` wires everything -- initialises repositories and services, creates the session provider, then maps CLI commands to controller constructors with their dependencies injected.

**Key patterns**:
- Controllers define small, consumer-side interfaces for the services they need (e.g. `projectLister`, `selecter`, `sessionLauncher`)
- Session provider abstraction (`SessionProvider` interface) supports tmux and VS Code, selectable via config or `--provider` flag
- `ProviderConfigOverride` wraps base config to override the session provider at runtime
- Generic `WithSpinner[T any]` function for async operations with terminal spinner

## Code Style

- **Formatting**: `gofmt` enforced by CI
- **Imports**: stdlib, then external, then internal (`github.com/danielronalds/projman/...`)
- **Interfaces**: Define small interfaces in the consumer package, not the provider
- **Naming**: camelCase unexported, PascalCase exported; receivers are single lowercase letter
- **Constructors**: `NewXxxType()` returning concrete type (not pointer for small structs)
- **Errors**: Return with context via `fmt.Errorf("context: %v", err.Error())`; write user-facing errors to stderr
- **Types**: Use type aliases for clarity (e.g. `type projectName = string`)
- **Tests**: Table-driven tests with standard `testing` package

## CI Pipeline

GitHub Actions runs a clint pipeline (`.github/pipelines/ci.yaml`) with steps: formatting -> vet -> static analysis -> tests -> dependency check.
