# Project Organization Rules

- **Reuse before inventing:** Always check for existing scripts, TODOs, and documentation before creating new scripts or utilities. Use and move existing scripts (e.g., `scripts/cube_server_control.sh`) to their appropriate locations (e.g., `cube-server/scripts/`) instead of duplicating functionality.
- **Follow TODOs:** Prioritize items in `TODOs.md` and similar files before introducing new solutions.
- **Centralize control scripts:** Place all server control and orchestration scripts under the relevant service directory (e.g., `cube-server/scripts/`).


# Build, Test, and Enhancement Action Items

## 1. Build and Test Orchestration
1. Move all end-to-end and integration test scripts to `testing/end2end/` and `testing/scripts/` for consistency and discoverability.
2. Use Makefile targets in `testing/Makefile` to orchestrate all e2e and integration tests, so CI and developers can run them with a single command.
3. Ensure all test scripts are self-contained, robust, and clean up after themselves (e.g., kill background servers, clean temp files).
4. Add clear output and error reporting to all test scripts for easier debugging.
5. Ensure all test scripts and Makefiles use `set -euo pipefail` and fail fast on errors.
6. Add logging and timestamps to test output for easier CI debugging.
7. Use environment variables for all test configuration (ports, credentials, etc.) to avoid hardcoding in scripts.
8. Add a test for "port already in use" and other common failure scenarios.
9. Add a test for missing or invalid credentials in simulation and direct modes.
10. Add a test for CLI help output and error messages for all major commands.
11. Add a test for Makefile targets themselves (e.g., check that all expected targets exist and work).

## 2. Test Coverage and Structure
1. Consider using a test runner (e.g., Bats, shunit2, or Go-based e2e test harness) for more structured shell test cases and assertions.
2. Add parallel test execution support for non-conflicting tests to speed up CI.
3. Add coverage reporting for e2e and integration tests if possible.
4. Add more integration tests (including edge cases and failure scenarios).
5. Expand test coverage for all simulation endpoints (Azure, AWS, GCP, validation, etc.), including edge cases and error handling.

## 3. Build System Optimization
1. Unify build logic in a single root Makefile that delegates to submodule Makefiles for build, test, and clean targets.
2. Remove duplicate or legacy build/test logic from submodules and scripts.
3. Use Makefile variables for binary locations, build flags, and environment setup to avoid hardcoding paths in scripts.
4. Add a `make lint` target at the root and in each module to enforce code quality before builds/tests.
5. Add a `make check` or `make ci` target that runs all lint, build, and test steps in the correct order for CI.

## 4. Developer Experience and Documentation
1. Document all Makefile targets and test scripts in the main README and/or ENHANCE.md for developer onboarding.
2. Add documentation and usage examples for advanced features.
3. Review and migrate actionable items from `TODOs.md` to this file for structured enhancement tracking.

## 5. Advanced Features and Future Work
1. Improved error handling and user feedback.
2. Unified provider abstraction for all object storage operations.
3. Support for additional S3-compatible providers.
4. CLI flags for advanced S3 features (versioning, lifecycle, ACLs, etc).
5. Better credential management and fallback logic.
6. Enhanced debug and logging controls (toggle verbosity).


# ENHANCE.md


---

## Advanced Hetzner S3 Operations with AWS SDK

Hetzner Object Storage is S3-compatible, so you can use the AWS S3 SDK (`github.com/aws/aws-sdk-go-v2/service/s3`) for advanced operations:

- Multipart uploads
- Object versioning
- Lifecycle rules
- Access control lists (ACLs)
- Presigned URLs
- Object tagging
- Server-side encryption
- Bucket policies
- Cross-region replication
- Event notifications

**Usage:**
- Set the endpoint to `https://<region>.your-objectstorage.com`.
- Use Hetzner S3 access/secret keys.
- All AWS S3 API features supported by Hetzner can be used directly.

---

## Next Steps for `mt` (from TODOs)

- Review and migrate actionable items from `TODOs.md` to this file for structured enhancement tracking.
- Prioritize:
  - Improved error handling and user feedback.
  - Unified provider abstraction for all object storage operations.
  - Add more integration tests (including edge cases and failure scenarios).
  - Support for additional S3-compatible providers.
  - CLI flags for advanced S3 features (versioning, lifecycle, ACLs, etc).
  - Better credential management and fallback logic.
  - Enhanced debug and logging controls (toggle verbosity).
  - Documentation and usage examples for advanced features.

---

## General Simulation Persistence

Enable general file-based simulation persistence for cube-server (see shared/providers/hetzner/objectstorage.go and other providers). Use the `CUBE_SERVER_SIM_PERSIST` environment variable to specify the persistence file path for test isolation. Test scripts should set this to a temp file and clean up after each run to ensure repeatable, isolated tests.

---

Add further enhancement ideas here as needed.
