# TODOs for punchbag-cube-testsuite

## Milestones
- [ ] Finalize example coverage and documentation for all scenarios
- [ ] Set up automated integration testing (CI, example validation)
- [ ] Expand documentation and developer onboarding materials
- [ ] Enhance generator/transformator for more providers/resources and config-driven workflows
- [ ] Prepare release process and distribution (versioning, changelog, binaries)
- [ ] Add StackIT, Hetzner, and IONOS object storage support (abstraction + mock logic) [medium priority]
- [ ] Enhance terraform-provider via punchbag generator (high priority):
    - Act as a multi-layer-cloud provider to run simulations against cube-server
    - Expose an interface for werfty-transformator to migrate Terraform code between cloud providers
    - Step-by-step suggestions:
        1. Refactor punchbag generator to expose a Go API and CLI for provider simulation and migration.
        2. Implement a Terraform provider plugin that proxies resource CRUD to punchbag/cube-server for simulation.
        3. Add a migration interface for werfty-transformator to convert Terraform HCL between providers using punchbag logic.
        4. Document provider mapping, simulation, and migration workflows.
        5. Add integration tests for simulation and migration via the Terraform provider.
- [ ] Add CLI support to multitool for managing object storage across all clouds (high priority):
    - Support direct management (create, list, delete, etc.) for AWS, Azure, GCP, StackIT, Hetzner, IONOS
    - Support management via simulation/proxy through sim-server or cube-server
    - Step-by-step suggestions:
        1. Extend multitool CLI to add unified object storage commands (create, list, delete, etc.)
        2. Implement provider selection and config (direct vs. proxy/simulation)
        3. Integrate with sim-server/cube-server APIs for simulation mode
        4. Add tests and documentation for all supported providers and modes
- [ ] Modularize monorepo for standalone multitool/mt releases and shared Go modules (medium priority):
    - Add go.mod to multitool/ and shared/ directories
    - Use go.work at repo root for local development
    - Ensure all imports use canonical module paths
    - Document and automate release process for mt binary and shared Go modules
    - Document the new modular/release workflow for multitool/mt and shared modules (medium priority)
    - Automate the release process for mt and shared modules (medium priority)
    - Set up GitHub Actions or CI for building and releasing mt and shared modules (medium priority)
    - Test the CLI and shared modules in both local development and release scenarios (medium priority)
    - Add CLI support to multitool/mt for Azure logging and Application Insights management (medium priority):
        - Implement logging and App Insights resource management for Azure first
        - Plan enhancements for other providers (AWS, GCP, etc.) in future steps

## Next Steps
- [ ] Finalize example coverage and documentation for all scenarios
- [ ] Set up automated integration testing (CI, example validation)
- [ ] Enhance generator/transformator for more providers/resources and config-driven workflows
- [ ] Expand documentation and developer onboarding materials
- [ ] Prepare release process and distribution (versioning, changelog, binaries)
- [ ] Fix and re-enable failing provider simulation API tests in `server/api/provider_simulation_test.go` (currently failing with 400/Invalid request body and nil interface conversion). Ignore for now, revisit after break.

## Low Priority / Future
- [ ] Integration workflow between generator and backend/server
- [ ] Add more advanced provider transformation features as needed
- [ ] Add support for StackIT Object Storage in generator
- [ ] Further automate the Terraform generation workflow (e.g., config-driven, batch generation, etc.)
- [ ] Add more provider pairs and conversion logic to `werfty-transformator`
- [ ] Update all source files at the root level to include an SPDX license comment for AGPL-3.0-only, and add a notice with the commit hash and author info.

## Notes
- There should be only one `multitool` directory in the workspace. Remove or merge any duplicates to avoid confusion and Go module issues.

---

**Note:**
- All new features and refactors should use the unified model layer in `shared/`.
- Keep Terraform provider code and backend logic strictly separated.
- Ensure all documentation, scripts, and instructions clearly distinguish between:
    - `multitool/` (the source directory)
    - `mt` (the CLI binary built from multitool)
- Always run the CLI as `./mt ...` from within the `multitool` directory, not as `./multitool/mt` or similar.
- Update all README, help, and onboarding docs to avoid confusion between the directory and the binary.
- Always check and set the correct working directory before running tests or commands. Use an environment variable (e.g., PUNCHBAG_BASE_DIR) as the base directory for all scripts and automation.

- Action plan for object storage provider integration:
  1. Ensure each provider uses the best/official SDK or API:
     - AWS: AWS SDK for S3
     - Azure: Azure SDK for Go (Blob Storage)
     - GCP: Google Cloud Storage SDK for Go
     - Hetzner: S3-compatible API (AWS SDK)
     - StackIT/IONOS: S3-compatible or official SDK
  2. Integrate and test each provider in the CLI, starting with Hetzner S3.
  3. Add full debug/error output for all create/list/delete operations.
  4. Only fall back to mock implementations if no credentials/config are found.
  5. Document and validate all endpoints, credentials, and region handling.

- Start with Hetzner S3: ensure real bucket creation, listing, and error handling work end-to-end.
