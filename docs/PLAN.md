# SAO TOC Server Implementation Plan

## Goal
Build a TOC server using `urfave/cli/v3` with:
- Root command: `sao`
- Subcommands: `server`, `ui`
- `ui` command serves a React app embedded in the Go binary
- Embedded NATS runtime for server messaging
- TAK Server feature parity as the target architecture
- CoT router as a core runtime component
- Config lifecycle at `/etc/sao/config.hcl`:
  - Create default config from embedded template when missing
  - Load config for runtime startup

## Scope
### In Scope
- Initialize Go module dependencies for `urfave/cli/v3`
- Implement CLI entrypoint and root command metadata
- Implement `server` command startup flow for TOC server
- Implement embedded NATS startup and shutdown wiring
- Implement TAK-compatible CoT ingestion, routing, and fanout pipeline
- Implement config create/load behavior at `/etc/sao/config.hcl`
- Store default config as an embedded file in the config package
- Build and embed React UI static assets for binary-only deployment
- Serve embedded React assets from the `ui` command
- Provide a dedicated Go client library for server APIs and CoT publish flows
- Add command help/usage text, startup validation, and basic error handling
- Document how to run and configure the system

## Architecture Principles
- Keep modules small, single-purpose, and named by responsibility.
- Keep command-layer code thin; move runtime logic into internal packages.
- Keep CoT in its own module with no transport-specific logic leaked into other packages.
- Use explicit interfaces at package boundaries (config, nats, toc, cot).
- Avoid cross-package coupling; depend on abstractions, not concrete internals.
- Prefer clear data flow over clever patterns so new contributors can trace behavior.
- Co-locate tests with each module and validate behavior at boundaries.
- Enforce standard tooling gates (`go fmt`, `go vet`, `go test`, `golangci-lint`).

### Out of Scope (Initial Milestone)
- 100% TAK ecosystem parity in the first release
- Advanced clustering, federation, authn/authz, or multi-node orchestration
- Advanced UI feature set beyond baseline operational views
- Production packaging/release automation

## Milestones

## Milestone 1: CLI Skeleton (`urfave/cli/v3`)
- Add dependency: `github.com/urfave/cli/v3`
- Create `main.go` with app setup and execution
- Define root command name as `sao`
- Register subcommands: `server`, `ui`
- Validate `sao --help` output

## Milestone 2: Config Bootstrap and Load
- Define config model for TOC server + NATS
- Add default config template file in config folder (HCL)
- Embed default config template using Go `embed`
- Implement `EnsureConfig(path)` behavior:
  - If `/etc/sao/config.hcl` is missing, create parent dir and write embedded template
  - If present, read and parse existing config
- Add validation and actionable startup errors
- Validate first-run (create) and subsequent run (load) flows

## Milestone 3: TOC Server + Embedded NATS Startup
- Implement `server` command action to:
  - Load config from `/etc/sao/config.hcl`
  - Start embedded NATS server with config-driven options
  - Start TOC server runtime components
- Add graceful shutdown (signals/context cancellation)
- Validate `sao server --help` and runtime startup path

## Milestone 4: CoT Router (TAK-Compatible Core)
- Implement CoT intake endpoints (protocols selected by config)
- Parse and validate CoT payloads
- Normalize envelope/metadata for internal routing
- Route CoT events by subscription and mission/channel rules
- Fanout routed CoT events to connected consumers
- Add delivery metrics and traceable message IDs for debugging
- Keep all CoT parsing/routing contracts in dedicated `internal/cot` module

## Milestone 5: Embedded React UI Command
- Create React app source under `web/ui`
- Build static assets to `web/ui/dist`
- Embed `web/ui/dist` into Go binary with `//go:embed`
- Implement `ui` command to serve embedded static assets over HTTP
- Support SPA fallback (serve `index.html` for non-asset routes)
- Add smoke tests for CLI, config bootstrap, CoT routing, and UI serving

## Milestone 6: Client Library (`pkg/client`)
- Create a reusable Go client package for SAO server APIs
- Implement client configuration (base URL, auth token, HTTP client injection)
- Add health check and CoT publish methods
- Add unit tests with `httptest` for request and error behavior
- Document basic client usage in project docs

## Milestone 7: Go Project Quality Baseline
- Add repository quality tooling (`Makefile`, lint config, editor config)
- Define contributor conventions in `CONTRIBUTING.md`
- Standardize local validation workflow (`make verify`, optional `make lint`)
- Keep CI-ready checks lightweight and deterministic

## Milestone 8: Repository Governance and Automation
- Add CI workflow for format, vet, test, race test, and lint
- Add Dependabot automation for Go modules and GitHub Actions
- Add pull request template to standardize change quality
- Add ownership and security policy documentation

## Proposed File Layout
```
.
├── main.go
├── pkg/
│   └── client/
│       ├── doc.go
│       ├── client.go
│       ├── options.go
│       ├── health.go
│       ├── cot.go
│       └── client_test.go
├── internal/
│   ├── cli/
│   │   ├── app.go
│   │   ├── command_server.go
│   │   └── command_ui.go
│   ├── config/
│   │   ├── default.hcl
│   │   ├── embed.go
│   │   ├── model.go
│   │   ├── load.go
│   │   └── defaults.go
│   ├── nats/
│   │   └── embedded.go
│   ├── toc/
│   │   ├── server.go
│   │   └── runtime.go
│   ├── ui/
│       ├── assets.go
│       └── server.go
│   └── cot/
│       ├── event.go
│       ├── parser.go
│       ├── router.go
│       └── subscriptions.go
├── web/
│   └── ui/
│       ├── src/
│       ├── public/
│       └── dist/
└── docs/
    └── PLAN.md
```

## Task Checklist
- [ ] Add `urfave/cli/v3` dependency
- [ ] Implement root command `sao`
- [ ] Register `server` and `ui` subcommands
- [ ] Add embedded default HCL config file under config package
- [ ] Implement config create/load for `/etc/sao/config.hcl`
- [ ] Implement embedded NATS startup/shutdown
- [ ] Implement TOC server startup path in `server` command
- [ ] Create dedicated `internal/cot` module boundaries and interfaces
- [ ] Implement TAK-compatible CoT parser and validation in `internal/cot`
- [ ] Implement CoT routing and fanout logic in `internal/cot`
- [ ] Add CoT routing observability (metrics/logging IDs)
- [ ] Create `pkg/client` reusable Go client library
- [ ] Add server health and CoT publish client methods
- [ ] Add unit tests for client request/response behavior
- [ ] Add baseline Go quality tooling and contributor standards
- [ ] Add CI and dependency automation (`.github/workflows`, Dependabot)
- [ ] Add PR template, CODEOWNERS, and security policy
- [ ] Add React UI source and build pipeline (`web/ui`)
- [ ] Embed built UI assets into the Go binary
- [ ] Implement `ui` command static server with SPA fallback behavior
- [ ] Enforce package boundaries and avoid circular dependencies
- [ ] Add focused unit tests per module (`config`, `nats`, `toc`, `cot`)
- [ ] Add help text, examples, and startup validation
- [ ] Verify startup, shutdown, and CoT routing execution paths
- [ ] Update project documentation with config details

## Risks and Mitigations
- Risk: Early command structure may not match future architecture  
  Mitigation: Keep command actions thin; delegate logic to internal packages.
- Risk: `/etc/sao/config.hcl` may require elevated permissions on some hosts  
  Mitigation: Return clear permission errors and document required privileges.
- Risk: Embedded NATS startup may fail due to port conflicts  
  Mitigation: Make ports configurable and include conflict diagnostics.
- Risk: TAK compatibility expectations can exceed initial implementation scope  
  Mitigation: Define a staged compatibility matrix and publish supported features per milestone.
- Risk: CoT routing bugs can cause message loss or loops  
  Mitigation: Add routing tests, dedup/loop-prevention rules, and delivery metrics.
- Risk: Inconsistent UX between subcommands  
  Mitigation: Standardize naming, descriptions, and error formatting.

## Definition of Done
- `sao --help` shows root command and both subcommands
- First `sao server` run creates `/etc/sao/config.hcl` when missing
- First `sao server` run writes config from embedded template in config package
- Subsequent `sao server` run loads `/etc/sao/config.hcl` and starts embedded NATS
- `sao server` ingests and routes CoT events through defined routing rules
- Core TAK-compatible CoT routing scenarios pass smoke/integration tests
- `sao ui` serves the embedded React application from the binary
- CoT implementation is isolated in its own module and reusable by runtime components
- `pkg/client` can connect, health check, and publish CoT with clear errors
- Core modules stay isolated and understandable (CLI, config, NATS, TOC, CoT, UI)
- Documentation reflects command usage, config path, and runtime intent

## Resolved Decisions
- Config path override is supported via `--config` and `SAO_CONFIG` env var.
- Embedded NATS defaults to local-only IPC boundary (`127.0.0.1:4222`) with configurable host/port.
- Readiness contract is `GET /v1/ready` with checks:
  - `config_loaded`
  - `nats_ready`
  - `runtime_ready`
- v1 mandatory TAK scope includes chat, data package, mission, and auth tracks.
- Initial CoT transports are TCP and TLS.
- React toolchain is Vite for `web/ui`.
