# TODOs for punchbag-cube-testsuite

## High Priority (Done)
- [x] Unify all provider/status/test models in `shared/` and update all components to use them
- [x] Remove duplicate/legacy types and fix all import paths
- [x] Clean up workspace: remove `cube-server` and `terraform-provider` directories
- [x] Ensure all proxy endpoints are in `server` and build/test cleanly
- [x] Integrate shared models into `werfty-generator`
- [x] Expand `werfty-generator` to support more Terraform resources/providers
- [x] Refactor and modernize `werfty-transformator` (modular, extensible, clean CLI)

## Next Steps
- [ ] Move all Terraform provider (HCL/DSL) generation logic to `werfty-generator`
- [ ] Keep all backend/proxy/simulation logic in `server`
- [ ] Use only unified types from `shared/` everywhere

## Low Priority / Future
- [ ] Integration workflow between generator and backend/server
- [ ] Add more advanced provider transformation features as needed
- [ ] Add support for StackIT Object Storage in generator
- [ ] Further automate the Terraform generation workflow (e.g., config-driven, batch generation, etc.)
- [ ] Add more provider pairs and conversion logic to `werfty-transformator`

---

**Note:**
- All new features and refactors should use the unified model layer in `shared/`.
- Keep Terraform provider code and backend logic strictly separated.
