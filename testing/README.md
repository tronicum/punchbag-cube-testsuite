# Testing Overview

This directory contains all integration and end-to-end tests for the punchbag-cube-testsuite project.

## Testing Concept

- **All orchestration and automation must use the multitool CLI (`./multitool/mt`) and server control scripts.**
- **No test or script may manipulate simulation state or server processes directly.**
- **All test scripts must be workspace-root-relative and robust to CI/local runs.**
- **All server startup/shutdown must use `scripts/cube_server_control.sh` for port management and isolation.**
- **All test scripts must check and set the correct working directory before running commands.**
- **All test orchestration must use Makefiles or documented scripts.**

## End-to-End Hetzner S3 Simulation Test

- The script `testing/end2end/end2end_hetzner_s3_sim.sh` runs a full simulation of Hetzner S3 object storage using only the multitool CLI and the simulation server.
- The test flow:
  1. Build all binaries using the root Makefile.
  2. Start the simulation server in simulate mode using the control script on a test port.
  3. Use the multitool CLI to create, list, and delete a test bucket via the simulation endpoint.
  4. Validate that bucket operations work as expected.
  5. Stop the simulation server and clean up all state.
- The test is fully automated and suitable for CI.

## Adding/Updating Tests

- Place new end-to-end tests in `testing/end2end/` and document them here.
- Always update this README when adding new test flows or orchestration patterns.
