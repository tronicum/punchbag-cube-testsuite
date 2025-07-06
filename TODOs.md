# TODOs for punchbag-cube-testsuite

## Milestones
- [ ] Finalize example coverage and documentation for all scenarios
- [ ] Set up automated integration testing (CI, example validation)
- [ ] Expand documentation and developer onboarding materials
- [ ] Enhance generator/transformator for more providers/resources and config-driven workflows
- [ ] Prepare release process and distribution (versioning, changelog, binaries)

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
