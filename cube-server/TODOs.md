# 2025-08-02: API Refactor and Endpoint Structure

- [x] Simulation endpoints are now under `/api/v1/simulate/`.
- [ ] Proxy endpoints to be implemented under `/api/v1/proxy/`.
- [ ] Direct endpoints to be implemented under `/api/v1/direct/` (if needed).
- [ ] Remove or deprecate any ambiguous or legacy endpoints (e.g., `/api/simulate/objectstorage/...`).
- [ ] Update CLI (multitool) to use the correct base path for each mode (`--simulate`, `--proxy`, `--direct`).
- [ ] Update all test scripts to use the new endpoint structure.
- [ ] Add deprecation warnings for legacy endpoints during transition.

See also: ENHANCE.md, ARCH.md for architectural rules and patterns.
# cube-server TODOs

## API Modes and Endpoint Structure
- [x] All simulation endpoints must be under `/api/v1/simulate/`.
- [x] All proxy (real provider) endpoints must be under `/api/v1/proxy/`.
- [x] All direct endpoints (if needed) must be under `/api/v1/direct/`.
- [x] CLI (multitool) must select the correct base path based on mode (`--simulate`, `--proxy`, `--direct`).
- [ ] Remove or deprecate any ambiguous or legacy endpoints (e.g., `/api/simulate/objectstorage/...`).
- [ ] Document all new endpoints and modes in ARCH.md and ENHANCE.md.

## Migration/Refactor
- [ ] Migrate any remaining simulation/proxy endpoints to the new structure.
- [ ] Ensure all new features follow this pattern.
- [ ] Add tests for all new endpoints.

## Notes
- Never mix simulation and proxy logic in a single handler or endpoint.
- All orchestration and automation must use the correct endpoint for the mode.
- Update CLI and test scripts to match the new endpoint structure.
- Add deprecation warnings for any legacy endpoints during transition.
