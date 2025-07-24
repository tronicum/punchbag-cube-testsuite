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
