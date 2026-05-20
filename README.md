# Deployment Events Service (Go + Gin)

Lightweight backend service for serving deployment event data.

## Requirements

- Go 1.25+ (required by Gin v1.12; `go.mod` uses Go 1.26.3)

### Install Go

Check if Go is already installed:

```bash
go version
```

If that fails, install Go using one of these options:

**macOS (Homebrew)**

```bash
brew install go
```

**Linux (Debian/Ubuntu)**

```bash
sudo apt update
sudo apt install -y golang-go
go version
```

If the distro package is older than 1.25, install from the official tarball instead: https://go.dev/dl/

**All platforms (official installer)**

1. Download the installer for your OS from https://go.dev/dl/
2. Install it
3. Restart your terminal
4. Verify:

```bash
go version
```

## Run locally (under 2 minutes)

```bash
make run
```

Or without Make:

```bash
go mod tidy
go run ./cmd/server
```

Server starts on `http://localhost:8080`.

Quick sanity checks:

```bash
curl "http://localhost:8080/health"
curl "http://localhost:8080/deployments"
curl "http://localhost:8080/deployments/deploy_001"
```

## API

### `GET /deployments`

List deployments. Supports optional query filters:

- `service` (example: `billing-api`)
- `status` (one of: `success`, `failed`, `cancelled`; example: `failed`)

Example:

```bash
curl "http://localhost:8080/deployments?service=billing-api&status=failed"
```

Response shape:

```json
{
  "data": [
    {
      "id": "deploy_001",
      "service": "billing-api",
      "status": "failed",
      "duration": 123,
      "timestamp": "2025-04-01T09:00:00Z",
      "commit_sha": "a00100"
    }
  ]
}
```

Invalid status response (`400 Bad Request`):

```json
{
  "error": {
    "code": "invalid_status",
    "message": "invalid status \"in_progress\"; valid values: success, failed, cancelled"
  }
}
```

### `GET /deployments/:id`

Fetch a single deployment by id.

Example:

```bash
curl "http://localhost:8080/deployments/deploy_001"
```

Success response (`200 OK`):

```json
{
  "data": {
    "id": "deploy_001",
    "service": "billing-api",
    "status": "failed",
    "duration": 123,
    "timestamp": "2025-04-01T09:00:00Z",
    "commit_sha": "a00100"
  }
}
```

Not found response (`404 Not Found`):

```json
{
  "error": {
    "code": "not_found",
    "message": "deployment not found"
  }
}
```

### `GET /health`

Simple health check endpoint.

Example:

```bash
curl "http://localhost:8080/health"
```

Response shape:

```json
{
  "status": "ok"
}
```

## Notes

- Data is seeded in-memory at startup with 50 mock deployment events.
- Seed data is randomly generated on each server start.
- Events span multiple services, statuses, durations, and timestamps.
- Default port is `8080` (override with `PORT` environment variable).

## Build, run, and test

```bash
make build
make run
make test
```