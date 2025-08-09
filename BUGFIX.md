# BUGFIX: Server Does Not Create Bucket in Simulation Mode

## Steps to Debug and Fix

1. **Confirm Server Startup and Logging**
   - Ensure the server is started with the correct binary: `./cube-server/cube-server`.
   - Check that Gin middleware for request logging is active.
   - Verify that debug output appears in the server logs.

2. **Check Simulation Handler**
   - Confirm the simulation handler for bucket creation is registered and matches the CLI request route.
   - Ensure the handler logs debug output when called.

3. **Verify Simulation Logic**
   - Inspect `shared/simulation/service.go` for bucket creation logic.
   - Confirm the handler calls the simulation service and that the logic is correct.

4. **Run End-to-End Test**
   - Use `testing/end2end/end2end_hetzner_s3_sim.sh` to run the test.
   - Ensure the server is running and accessible at the expected port.
   - Check logs for request/response flow and debug output.

5. **Workspace Rules**
   - Always run/build/test the server from inside the `cube-server/` directory.
   - Always check `pwd` before running Go commands.
   - All scripts and binaries must use workspace-root-relative paths.

6. **If Bucket Is Not Created**
   - If the handler is not called, verify the route and request.
   - If the handler is called but the bucket is not created, debug the simulation logic.
   - Add or review debug output as needed.

---

_Always stop and resolve errors before proceeding to the next step. Follow workspace rules strictly._
