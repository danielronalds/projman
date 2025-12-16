# AGENTS.md

## Project Overview
projman is a CLI tool for managing development projects. It integrates with tmux to launch sessions and supports local project discovery, GitHub remote repos, and project templating.

## Build & Test Commands
- Build: `go build ./...`
- Test all: `go test ./...`
- Test single package: `go test ./internal/services`
- Test single test: `go test ./internal/services -run TestFunctionName`
- Format check: `gofmt -l .` (CI enforces formatting)
- Format fix: `gofmt -w .`

## Code Style
- **Formatting**: Use `gofmt` - CI will fail on unformatted code
- **Imports**: Group stdlib first, then external packages, then internal (`github.com/danielronalds/projman/...`)
- **Architecture**: Controllers handle CLI args, Services contain business logic, Repositories manage config/data
- **Interfaces**: Define small interfaces in the consumer package (see controller files for examples)
- **Naming**: Use camelCase for unexported, PascalCase for exported; receivers are single lowercase letter
- **Errors**: Return errors with context using `fmt.Errorf("context: %v", err.Error())`; write to stderr with `fmt.Fprintf(os.Stderr, ...)`
- **Constructors**: Name as `NewXxxType()`, return concrete type (not pointer for small structs)
- **Types**: Use type aliases for clarity (e.g., `type projectName = string`)
