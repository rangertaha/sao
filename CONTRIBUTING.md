# Contributing

## Development Requirements
- Go 1.22+
- `golangci-lint` installed locally

## Project Layout Guidelines
- Put executables in `cmd/<name>`.
- Keep non-exported server internals in `internal/...`.
- Keep reusable libraries for external consumers in `pkg/...`.
- Keep modules small and focused by responsibility.

## Coding Conventions
- Use the Go standard library first when possible.
- Return wrapped errors with context (for example: `fmt.Errorf("load config: %w", err)`).
- Keep functions short and clear; extract helpers when logic grows.
- Prefer constructor + option patterns for reusable packages.
- Avoid package-level mutable state unless strictly required.

## Quality Gates
Run these before opening a PR:

```sh
make verify
```

If `golangci-lint` is available, also run:

```sh
make lint
```

## Testing
- Co-locate tests with the package under test.
- Use `httptest` for HTTP boundaries.
- Favor table-driven tests for parser and router logic.
