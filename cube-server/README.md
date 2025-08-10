# Cube Server Documentation

This document covers usage, configuration, and architecture for the cube-server component.

## Overview
- The cube-server provides unified orchestration, simulation, and API endpoints for all supported cloud providers.
- All server logic is located in the `cube-server/` directory.

## Usage
- Build: `make build` (from cube-server directory)
- Start: `make start` or run the built binary
- Test: `make go-tests`

## Endpoints
- See `api/openapi.yaml` for full API specification.
- Simulation endpoints for Azure, AWS, GCP, Hetzner, etc.

## Debug Mode

To start the server in debug mode (verbose logging, error details), use the `--debug` flag:

```
go run main.go --debug
```

Or with a custom port:

```
go run main.go --port 9090 --debug
```

This enables development logging and prints debug information to the console.

## Developer Notes
- All provider logic must use the shared/ library abstraction.
- No direct provider/model code outside shared/.
- Server orchestration logic is modular and maintainable.

## References
- See main README.md for high-level project info and links to other app docs.
- See `shared/README.md` for shared library API and usage.
- See `multitool/README.md` for CLI documentation.

---
SPDX-License-Identifier: AGPL-3.0-only
