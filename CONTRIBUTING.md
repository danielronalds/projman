# Contributing

This project uses [mise](https://mise.jdx.dev) for tool management and task running.

## Setup

```console
mise install
```

## Tasks

| Task | Description |
|------|-------------|
| `mise dev` | Run the project |
| `mise build` | Build the binary |
| `mise build:install` | Build and install to `~/.local/bin/` |
| `mise test` | Run all tests |
| `mise lint:fmt` | Check for unformatted code |
| `mise lint:vet` | Run go vet |
| `mise lint:static` | Run golangci-lint |
| `mise lint:deps` | Check for unused or missing dependencies |
| `mise clean` | Remove built binary |

## CI

A GitHub Actions pipeline runs on pushes to `main` and on pull requests. It runs formatting, vet, static analysis, tests, and dependency checks.
