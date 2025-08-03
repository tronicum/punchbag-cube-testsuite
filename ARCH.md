# Architecture Note: Generic AWS S3 Simulation Endpoint
- The system must provide a generic AWS S3 simulation endpoint under /api/v1/simulate/aws-s3/* to support SDK-compatible testing and integration.
- This endpoint should be provider-agnostic and simulate standard AWS S3 API behavior for all basic bucket/object operations.
# punchbag-cube-testsuite Architecture

## Component Overview

- **shared/**: Central Go module for all cloud/provider abstractions, models, and shared logic. All applications must use this for provider operations.
- **cube-server/**: Unified server and simulation component. Contains all simulation, API, and backend logic. Uses only the shared module for provider abstractions. All simulation and server logic is consolidated here.
- **sim/** and **sim-server/**: Legacy simulation modules. All unique logic should be migrated to cube-server. These will be deprecated and removed after migration.
- **multitool/**: Unified CLI tool. All commands (including k8sctl, k8s-manage) use shared abstractions and models. No direct provider logic outside shared.
- **werfty/**, **werfty-generator/**, **werfty-transformator/**: Modular applications for resource management, code generation, and transformation. Each is a separate Go module for maintainability.

## Architectural Rules

- All provider/cloud logic must reside in shared/.
- No direct provider logic or models outside shared/.
- All applications (cube-server, multitool, werfty, etc.) must use shared/ for cloud/resource operations.
- CLI structure: multitool provides top-level subcommands (k8sctl, k8s-manage, etc.) that are provider-agnostic and extensible.
- Simulation and server logic is unified in cube-server. Legacy sim/sim-server will be removed.
- All documentation, scripts, and usage must reference the multitool binary as ./multitool/mt.

## Migration Plan (for sim/sim-server/cube-server)

1. [x] Move all simulation logic into cube-server.
2. [x] Remove sim/ and sim-server/ after migration.
3. [x] Ensure cube-server uses only shared/ for provider logic.
4. [x] Update go.work and all references to use cube-server for simulation and backend operations.
5. [x] General simulation persistence: The cube-server now supports file-based persistence for simulation state (e.g., buckets, objects) for all providers. When the environment variable CUBE_SERVER_SIM_PERSIST is set, all simulation state is loaded from and saved to the specified file (default: /tmp/cube_server_sim_buckets.json). This enables robust, repeatable end-to-end tests and allows the simulation server to be restarted without losing state. Test scripts should set/clean this file for isolation.

## Next Steps (Post-Migration)

1. Fix simulation handler logic to ensure correct status codes and endpoint behavior.
2. Expand test coverage for all simulation endpoints (Azure, AWS, GCP, validation, etc.), including edge cases and error handling.
3. Audit all server and CLI code to enforce usage of shared/ for provider logic; remove any direct provider/model code outside shared/.
4. Update README.md and developer docs with new endpoint details, usage examples, and architectural rules.
5. Add CI checks to enforce module hygiene, run tests, and validate shared usage.
6. Review TODOs.md for remaining migration, refactor, and feature tasks; prioritize next CLI, provider, or integration features.


## Notes
- See TODOs.md for current sprint tasks and feature roadmap.
- See README.md for CLI usage and developer documentation.
- This file is the single source of truth for architecture and module boundaries.

## OS/Package Manager Testing Strategy

To ensure robust cross-platform support for OS detection and package management features, the following approach is used:

- All CLI binaries are always named `mt` for all Linux/arch builds.
- OS and package manager detection is tested using Docker containers for each major Linux distribution.
- The containers are only used for testing the package manager and OS detection features, not for general CLI use.
- All Dockerfiles for testing are located in `testing/docker/` and named by OS (e.g., `Dockerfile.ubuntu`, `Dockerfile.alpine`, etc.).
- Each Dockerfile copies the pre-built `mt` binary into the container and sets it as the entrypoint.
- Example test commands:
  - `docker build -f testing/docker/Dockerfile.ubuntu -t mt-ubuntu .`
  - `docker run --rm mt-ubuntu os-detect`
  - `docker run --rm mt-ubuntu list-packages`
- This ensures that the `mt` binary's OS/package manager features work as expected on all supported distributions.
- The containers are not intended for end-user use, only for CI/maintainer testing.

See `TESTING.md` for more details and usage instructions.

## Simulation Test Flow (Hetzner S3)
- All bucket operations (create, list, delete) must be performed via multitool CLI.
- No direct file or state manipulation is allowed in test scripts.
- Dummy S3 credentials can be injected via SIMULATE_DUMMY_S3_CREDS.
- [Planned] Dummy S3 buckets can be injected via SIMULATE_DUMMY_S3_BUCKETS (to be implemented).
- The simulation server persists state via CUBE_SERVER_SIM_PERSIST.

## Test Orchestration Rule
- All test setup, execution, and teardown must use multitool CLI commands only.
- Scripts may only orchestrate CLI calls and may not manipulate simulation state directly.
