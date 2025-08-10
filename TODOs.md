# TODO: Add generic AWS S3 simulation endpoint for SDK compatibility
- Implement a generic endpoint under /api/v1/simulate/aws-s3/* for simulating AWS S3 SDK-compatible operations (provider: generic-aws-s3)
- Ensure this endpoint is documented and tested
## TODO: Simulation/Mock Migration

- [x] Migrate all Hetzner S3 mock logic from multitool and shared/providers/hetzner to sim-server within cube-server.
- [x] Remove legacy mock server code from multitool.
- [x] Ensure all simulation endpoints are served via cube-server/sim-server abstraction.
- [x] Add general file-based simulation persistence to cube-server (shared/providers/hetzner/objectstorage.go and other providers). When CUBE_SERVER_SIM_PERSIST is set, all simulation state (e.g., buckets) is loaded/saved to a JSON file (default: /tmp/cube_server_sim_buckets.json). This enables robust, repeatable simulation tests and server restarts.
    - [x] Add env var/flag for persistence file path (CUBE_SERVER_SIM_PERSIST).
    - [x] Document in ARCH.md and ENHANCE.md.
    - [x] Ensure test scripts set/clean this file for isolation.
- [ ] Validate the migration by running all server-related tests and operations from cube-server/ only.

## Refactor Plan: Merge server and cube-server

- [ ] Audit both directories for unique logic, scripts, and configuration.
- [ ] Move all essential orchestration, simulation, and server logic into cube-server/.
- [ ] Remove any duplicate, legacy, or unused files from server/.
- [ ] Update all Makefiles, scripts, and documentation to reference only cube-server/ for server operations.
- [ ] Ensure all tests, CI, and developer workflows use the unified cube-server/.

## Modular Go Test Orchestration Rule

- [ ] Only define a `go-tests` target in a module Makefile if Go code and tests are present.
- [ ] The main orchestration in `testing/Makefile` should only call `go-tests` for modules with actual Go tests.
- [ ] Do not use shell logic to check for Go files; simply omit the target if there are no tests.
# Outstanding Test Failures (to fix later)


## Multitool/Cube-Server Modularization TODOs

*Moved to `multitool/TODOs.md` for better organization.*


- [ ] **generator**: Test failures due to mixed package names (`main` and `generator`) in the same directory, and missing internal package references. Example errors:
    - found packages main (aks.go) and generator (end2end_test.go) in generator/
    - package punchbag-cube-testsuite/generator/internal/generator is not in std

- [ ] **werfty**: Test failures due to missing/misconfigured imports from client and sharedmodels, plus struct field/type mismatches. Example errors:
    - package punchbag-cube-testsuite/client/cmd is not in std
    - unknown field CloudProvider in struct literal of type models.Cluster
    - undefined: sharedmodels.AKSTestRequest
    - cannot use cluster.Status (variable of string type models.ClusterStatus) as string value in struct literal

- [ ] Fix generator and werfty test errors after modular orchestration is validated and other priorities are complete.
## Documentation Modularization

- [ ] Move multitool documentation to multitool/README.md and reference it from the main README.md
- [ ] Keep main README.md compact, with only high-level project info and links to per-app docs
- [ ] Repeat for other applications (e.g., cube-server, k8sctl) as needed
# Medium Priority
- [ ] Unified Kubernetes (k8s) Management for All Clouds (direct mode)
    - [x] Scaffold direct mode k8s management for all supported providers (Azure, AWS, GCP, Hetzner, IONOS, StackIT, OVH, etc.)
        - [x] Add stub functions and CLI commands for create, update, scale, upgrade, and delete clusters and node pools
        - [x] Use a unified resource model and provider-agnostic CLI structure
        - [x] Add provider-specific stub files (e.g., `shared/providers/aws/k8s.go`, `shared/providers/gcp/k8s.go`, etc.)
        - [x] Document required API tokens/credentials for each provider
    - [ ] (Testing and implementation to follow once API tokens are available)
#
# SPDX-License-Identifier: AGPL-3.0-only
#
# Copyright (C) 2023-2025 tronicum@user.github.com
#
# This file is part of punchbag-cube-testsuite.
#
# punchbag-cube-testsuite is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# punchbag-cube-testsuite is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with punchbag-cube-testsuite.  If not, see <https://www.gnu.org/licenses/>.
# OVHcloud Integration
- [ ] Add OVHcloud as a supported cloud provider (shared/, multitool, etc.)
- [ ] Scaffold S3-compatible object storage service for OVHcloud (API, mocks, CLI integration)
- [ ] Integrate OVHcloud S3 logic using https://registry.terraform.io/providers/ovh/ovh/latest for Terraform-based operations and transformations
- [ ] When transforming or generating Terraform, use the OVHcloud provider as reference for OVH S3 resources (future: support more providers as needed)

## Low Priority

- [ ] Refactor the client/ module for maintainability, shared code usage, and API consistency (Low priority)
- [ ] Investigate enabling Copilot (AI assistant) to access GitHub for direct commits and merge requests. (Low priority)
## Notes on end2end_multitool_test.sh

- The variable `SIMULATOR_BIN` in `testing/end2end/end2end_multitool_test.sh` points to the multitool CLI binary. It is used to run multitool commands directly (such as objectstorage and resource simulations) and is not a client for the sim-server. This ensures all CLI calls use the correct binary location.

## Notes on Hetzner S3 Bucket Metadata
- [ ] Hetzner Object Storage (and all S3-compatible APIs) do not provide bucket creation or update timestamps via the S3 API. This is a limitation of the protocol and not the implementation. If richer metadata is needed, monitor Hetzner's hcloud API for future support.

## Storage Provider Abstraction Rule
- [ ] Enforce that *all* storage provider logic (not just in multitool) must use the shared/ abstraction. No direct provider/model code outside shared/.
# TODOs for punchbag-cube-testsuite
## Planned: Multitool kubeconfig-like config system

- [ ] Implement `multitool/.mtconfig/<profile>/config.yaml` profile system for multitool:
    - Each profile is a directory under `multitool/.mtconfig/`, e.g. `multitool/.mtconfig/aws-dev/config.yaml`.
    - Config file in YAML format, supports provider, region, credentials, endpoints, and custom settings.
    - Profile switching via CLI flag `--profile` or environment variable.
    - Provider-specific settings and secrets (static, ENV, or external secret managers).
    - Extensible for new fields/providers.

### Implementation Steps
1. Scaffold `multitool/.mtconfig/<profile>/config.yaml` loader in multitool.
2. Add CLI flag `--profile` and env var support.
3. Update provider selection logic to use config values.
4. Document usage and migration in README.md.
> **NOTE:**
> The *only* supported CLI binary is `./multitool/mt`.
> All documentation, scripts, and usage must reference this binary.
> Do **not** use `./mt`, `multitool/multitool`, or any other binary name/location.
> This is enforced in the Makefile and build process.

## Next Steps (Post-Migration)

- [ ] Fix simulation handler logic to ensure correct status codes and endpoint behavior.
- [ ] Expand test coverage for all simulation endpoints (Azure, AWS, GCP, validation, etc.), including edge cases and error handling.
- [ ] Audit all server and CLI code to enforce usage of shared/ for provider logic; remove any direct provider/model code outside shared/.
- [ ] Update README.md and developer docs with new endpoint details, usage examples, and architectural rules.
- [ ] Add CI checks to enforce module hygiene, run tests, and validate shared usage.
- [ ] Review TODOs.md for remaining migration, refactor, and feature tasks; prioritize next CLI, provider, or integration features.

> **ARCHITECTURE NOTE:**
> All cloud/provider abstractions must reside in the shared/ library. Application-specific logic (for multitool, werfty, punchbag server, etc.) should only adapt or extend the shared abstraction as needed for their context.
> All components (punchbag server, mt in proxy mode, mt in direct mode, werfty, etc.) must use the same shared abstraction layer for all cloud and resource operations. No direct provider logic or models should exist outside shared/.

## CLI Command Structure Rules (K8s Management)

- All Kubernetes-related CLI operations must be split into two distinct subcommands:
  - `k8sctl`: For kubectl-like operations (apply, get, exec, logs, etc.)
  - `k8s-manage`: For cluster lifecycle management (create cluster, delete cluster, scale, upgrade, etc.)
- Both subcommands must be at the top level of the CLI (i.e., `mt k8sctl ...` and `mt k8s-manage ...`).
- The `k8s-manage` subcommand must use a provider-agnostic structure, e.g.:
    - `mt k8s-manage create cluster --provider hetzner ...`
    - `mt k8s-manage delete cluster --provider aws ...`
    - `mt k8s-manage scale cluster ...`
    - (etc. for all supported providers)
- The `k8sctl` subcommand should mirror the UX of `kubectl` as closely as possible, but always operate through the multitool abstraction (never call kubectl directly unless explicitly requested).
- All new CLI features must be designed for maximum flexibility and extensibility, to support future providers and new resource types without breaking changes.
- Document these rules in all relevant developer docs and keep them up to date as the CLI evolves.


## Current Sprint Tasks (This Week)

### Phase A: Shared Library Integration (ACTIVE)
- [x] Move Azure commands from multitool to shared/ library
- [x] Create shared/providers/azure/ package with all Azure operations
- [x] Fix compilation issues in multitool azure.go (FIXED - cleaned up orphaned code)
- [x] Update multitool to use shared/providers/azure instead of local commands
- [x] Create shared/export package for data exchange
- [x] Standardize shared/ library API for cloud providers
- [x] Add shared/import packages for data exchange
    - [ ] Test all applications using shared/ library
        - [ ] Add all relevant modules to go.work
        - [ ] Ensure all submodules have correct go.mod files and import paths
        - [ ] Run go mod tidy in each module directory
        - [ ] Re-run all tests and verify shared library integration

### Phase B: Strategic Architecture (AFTER A)
- [x] Scaffold terraform-multicloud-provider structure
<!-- Completed tasks removed for clarity -->

## Milestones

<!-- Completed tasks removed for clarity -->
### IMMEDIATE: Shared Library Migration (CRITICAL PRIORITY)

- [x] Migrate Azure cloud operations to `shared/providers/azure/`:
<!-- Completed tasks removed for clarity -->
- [x] Create `shared/export` and `shared/import` packages:
    - Export cloud state to JSON/YAML (for werfty-generator)
    - Import configurations and validate (for werfty-transformator)
- [ ] Standardize shared/ library interfaces:
    - [ ] Common Provider interface for all clouds (review for consistency)
    - [ ] Standardized authentication and configuration (review/complete)
    - [ ] Consistent simulation vs direct mode handling (review)
- [ ] Complete shared module integration across all components:
    - [x] All applications use unified models from shared/
<!-- Completed tasks removed for clarity -->
    - [ ] Standardize API interfaces between multitool, generators, and cube-server
    - [x] Complete migration of common code to shared/ module
#### Step-by-step suggestions:
1. Review and finalize Provider interface in `shared/providers/interface.go`.
2. Standardize authentication/config patterns across all providers.
3. Refactor error handling and logging to use shared/log and shared/errors everywhere.
5. Update multitool, werfty-generator, and werfty-transformator to use only shared/ code for all cloud/model logic.

---

**TODO: Refactor multitool command structure for clarity and maintainability**

- [ ] Restructure multitool CLI commands to the following structure:
    - `mt aws ...`
        - `mt aws cloudformation ...` (all CloudFormation commands as subcommands)
        - `mt aws s3 ...` (all S3 commands as subcommands)
        - (future: `mt aws eks ...`, etc.)
    - `mt gcp ...` (all GCP-specific commands as subcommands)
    - `mt hetzner ...` (all Hetzner-specific commands as subcommands)
        - `mt docker container ...` (container management)
        - `mt docker image ...` (image management)
        - `mt docker registry ...` (registry login/logout/list, etc.)
    - `mt k8s ...` (provider-agnostic Kubernetes commands, with `--provider` flag)
    - `mt local ...`
        - `mt local os ...` (OS detection/info)
        - `mt local package ...` (package management)
        - `mt local file ...` (local file operations, validation, etc.)
    - `mt config ...` (global config management)
    - `mt test ...` (testing utilities)
    - `mt scaffold ...` (project scaffolding, if needed)

- [ ] Review and update documentation and help output to reflect the new structure
- [ ] Ensure backward compatibility or provide migration notes for users
    1. [x] Audit current shared/ usage across multitool, werfty-generator, werfty-transformator
    2. [x] Move remaining common models and utilities to shared/
    3. [x] Update import paths in all applications to use shared/
    4. [ ] Standardize error handling and logging across components
    5. [ ] Document shared module API and usage patterns


1. Standardize error handling and logging across all shared/ and consuming apps (multitool, werfty-generator, werfty-transformator, cube-server).
2. Review and unify the Provider interface and authentication/config patterns in shared/providers/interface.go and all provider implementations.
3. Document the shared module API and usage patterns for all teams.
4. Test all applications (multitool, werfty-generator, werfty-transformator) independently using only the shared/ library (no local model/util code).
5. (Optional) Add more provider-agnostic tests and CI checks to enforce shared/ usage and interface compliance.
- [ ] Add Azure DevOps support to multitool framework (medium priority):
    - Implement Azure DevOps project, pipeline, and repository management
    - Support for Azure Boards, Repos, Artifacts, and Test Plans operations
    - Integration with existing Azure provider functions and CI/CD workflows
    - Step-by-step suggestions:
        1. Add Azure DevOps provider commands to multitool CLI (projects, pipelines, repos, boards)
        2. Implement Azure DevOps REST API integration for real operations

- [ ] Enhance multitool CLI capabilities (high priority)

# Phase 4: Advanced Integration & CI
- [ ] Test all applications (multitool, werfty-generator, werfty-transformator) independently using only the shared/ library (no local model/util code)
- [ ] Add more provider-agnostic tests and CI checks to enforce shared/ usage and interface compliance
- [ ] Integration workflow between generator and backend/server
- [ ] Add more advanced provider transformation features as needed
- [ ] Add support for StackIT Object Storage in generator
- [ ] Further automate the Terraform generation workflow (e.g., config-driven, batch generation, etc.)
- [ ] Add more provider pairs and conversion logic to `werfty-transformator`
- [ ] Update all source files at the root level to include an SPDX license comment for AGPL-3.0-only

# Phase 5: New Features & Provider Support
#
# ---
#
## Next Steps for Kubernetes CLI Management (Status: July 2025)

1. **Adopt new CLI structure:**
    - [x] Add `k8sctl` as a top-level subcommand for kubectl-like operations.
    - [x] Add `k8s-manage` as a top-level subcommand for cluster lifecycle management (create, delete, scale, upgrade, etc.).
    - [x] Ensure both subcommands are provider-agnostic and easily extensible.
2. **Scaffold k8s-manage commands:**
    - [x] Implement `create cluster`, `delete cluster`, `scale cluster`, etc., for Hetzner and all other supported providers.
    - [x] Use a unified resource model and shared logic for all providers.
3. **Scaffold k8sctl commands:**
    - [x] Implement kubectl-like commands (get, apply, exec, logs, etc.) that work through the multitool abstraction.
    - [x] Ensure kubeconfig management is seamless for all providers.
4. **Update developer documentation:**
    - [x] Document the new CLI structure and rules for all contributors.
    - [x] Add usage examples for both `k8sctl` and `k8s-manage` (see README.md).
5. **Design for flexibility:**
    - [x] All new CLI features are designed for maximum flexibility and future extensibility.
    - [x] Avoid provider-specific hacks; always use shared abstractions and models.
6. **(Optional) Migrate existing k8s logic:**
    - [x] Refactor any existing k8s commands in multitool to fit the new structure.
    - [x] Deprecate or migrate old commands as needed.
- [ ] Add Azure DevOps support to multitool framework (medium priority)
- [ ] Add StackIT, Hetzner, and IONOS object storage support (abstraction + mock logic)

# Ongoing/Low Priority
- [ ] Investigate enabling Copilot (AI assistant) to access GitHub for direct commits and merge requests. (Low priority)
- [ ] Hetzner Object Storage (and all S3-compatible APIs) do not provide bucket creation or update timestamps via the S3 API. This is a limitation of the protocol and not the implementation. If richer metadata is needed, monitor Hetzner's hcloud API for future support.
        2. Use shared/export to generate data for other applications

## FIX PLAN: Modular Test Orchestration & Makefile Setup (July 2025)

- [ ] Audit all Makefiles in root and subdirectories for duplicate/conflicting targets and .PHONY declarations.
- [ ] Ensure each subdirectory (multitool, cube-server, generator, werfty) has a Makefile with build, clean, and test targets.
- [ ] Rewrite `testing/Makefile`:
    - Single .PHONY declaration at the top.
    - Modular test targets for each module: test-multitool, test-shared, test-cube-server, test-generator, test-werfty.
    - Main test target calls all modular test targets.
    - Remove legacy/duplicate targets and any use of `go test ../...`.
- [ ] Update root Makefile to only call `$(MAKE) -C testing test` for the test target.
- [ ] Validate by running `make test` from the root and ensure all module tests are executed separately.
- [ ] Document the modular test orchestration in developer docs and README.md.

### Step-by-step Fix Process
1. Audit and clean up all Makefiles for duplicate/conflicting targets and .PHONY declarations.
2. Ensure all subdirectory Makefiles exist and are correct.
3. Rewrite `testing/Makefile` for modular test orchestration.
4. Update root Makefile to delegate test logic to `testing/Makefile`.
5. Validate with `make test` and individual module test runs.
6. Document the setup and update developer notes.

## ATOMIC REFACTOR PLAN: Makefile Modularization & Test Delegation (July 2025)

- [ ] Refactor Makefile structure for maintainability, modularity, and future extensibility:
    - Root Makefile: Only orchestrates build, test, and clean. Delegates all test logic to `testing/Makefile` and server orchestration to application Makefiles (e.g., `cube-server/Makefile`).
    - `testing/Makefile`: Contains all test logic, including Docker-based, multiarch, cross-distro CLI tests. Flexible targets for future test types (e.g., `mt-tests-docker`). All test scripts moved to `testing/scripts/`.
    - `cube-server/Makefile`: Contains all server orchestration logic (start/stop/build/clean).
    - Avoid duplicate/conflicting targets across Makefiles. Ensure clear boundaries and documentation.
    - Document all changes and update developer docs.

### Sprint Steps
1. Write this refactor plan to TODOs.md (done).
2. Confirm with user (done).
3. Refactor root Makefile to delegate all test logic and server orchestration.
4. Move all test logic to `testing/Makefile`.
5. Move server orchestration logic to `cube-server/Makefile`.
6. Clean up legacy/duplicate targets and scripts.
7. Update documentation and developer notes.
8. Validate with test runs and CI.

- [ ] Move Hetzner S3 simulation mock (`NewHetznerS3Mock` and related code) from `multitool/pkg/client/hetzner_s3_mock.go` to `shared/providers/hetzner/objectstorage.go`.
- [ ] Move the CLI command `simulate-hetzner-s3` from `multitool/cmd/sim_hetzner_s3.go` to `cube-server/cmd/sim_hetzner_s3.go`.
- [ ] Register the simulation command in cube-server only, and update documentation to reflect the new location.
- [ ] Remove legacy simulation code from multitool after migration.
