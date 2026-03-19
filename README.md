# sao
Spatial Awareness Operator is a TAK compatible open source server

**⚠️ Under development: do not use**





The `ui` command serves a React web app embedded in the SAO binary.

## Commands

- `sao server` runs the TOC server, starts embedded NATS, and serves API endpoints.
- `sao ui` serves embedded UI assets over HTTP.

By default, configuration is loaded from `/etc/sao/config.hcl`. If the file does not
exist, SAO creates it from the embedded default template.
Set `SAO_CONFIG` to override the config path for local/non-root development.

Server probes:
- `GET /v1/health` basic process health
- `GET /v1/ready` readiness checks (`config_loaded`, `nats_ready`, `runtime_ready`)

## Development Standards

This repository follows modern Go project conventions:
- `cmd/` for binaries
- `internal/` for non-exported server internals
- `pkg/` for reusable public libraries (for example `pkg/client`)
- CI on push/PR via GitHub Actions
- Weekly dependency update automation via Dependabot

Recommended local checks:

```sh
make verify
```

Optional lint (when installed):

```sh
make lint
```

Project policies:
- Contribution guide: `CONTRIBUTING.md`
- Security policy: `SECURITY.md`

## Go Client Library

The repository now includes a reusable SAO client library at `pkg/client`.

Example:

```go
c, err := client.New("http://127.0.0.1:8080", client.WithBearerToken("token"))
if err != nil {
    // handle error
}

health, err := c.Health(ctx)
if err != nil {
    // handle error
}

_ = health
```

## Local Development

```sh
make run-server
make run-ui
```


