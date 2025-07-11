# TODOs for punchbag-cube-testsuite

## Project Architecture Overview

**Core Components:**
- **multitool/mt**: Multi-cloud CLI for direct cloud operations (can proxy via cube-server)
- **werfty-generator**: Creates new Terraform code from cloud state (JSON/YAML → .tf)
- **werfty-transformator**: Imports/exports Terraform code between cloud providers (.tf ↔ .tf)
- **terraform-multicloud-provider**: Generic Terraform provider for simulation and multi-cloud workflows
- **cube-server**: Simulation backend for testing and development
- **shared/**: Common code shared across all applications (models, utilities, sim-server integration)
- **client/**: API client and formatting utilities

## Current Sprint Tasks (This Week)

### Phase B: Strategic Architecture (ACTIVE)
- [x] Scaffold terraform-multicloud-provider structure
- [ ] Design provider plugin architecture using shared/ models

### Phase A: Fix & Test Immediately (NEXT)
- [ ] Fix compilation issues in client/pkg/output/formatter.go
- [ ] Add missing AWS dependencies to multitool
- [ ] Fix vendor directory inconsistencies
- [ ] Test multitool Azure functions in both direct and simulation modes
- [ ] Implement missing API client methods for cube-server integration

### Phase C: Expand Functionality (AFTER A)
- [ ] Add more cloud providers to multitool (AWS, GCP basics)
- [ ] Enhance werfty-generator capabilities
- [ ] Implement object storage commands across providers

## Milestones
- [ ] Finalize example coverage and documentation for all scenarios
- [ ] Set up automated integration testing (CI, example validation)
- [ ] Expand documentation and developer onboarding materials
- [ ] Prepare release process and distribution (versioning, changelog, binaries)
- [ ] Complete shared module integration across all components (high priority):
    - Ensure all applications use unified models from shared/
    - Standardize API interfaces between multitool, generators, and cube-server
    - Complete migration of common code to shared/ module
    - Step-by-step suggestions:
        1. Audit current shared/ usage across multitool, werfty-generator, werfty-transformator
        2. Move remaining common models and utilities to shared/
        3. Update import paths in all applications to use shared/
        4. Standardize error handling and logging across components
        5. Document shared module API and usage patterns
- [ ] Add StackIT, Hetzner, and IONOS object storage support (abstraction + mock logic) [medium priority]
- [ ] Add Azure DevOps support to multitool framework (medium priority):
    - Implement Azure DevOps project, pipeline, and repository management
    - Support for Azure Boards, Repos, Artifacts, and Test Plans operations
    - Integration with existing Azure provider functions and CI/CD workflows
    - Step-by-step suggestions:
        1. Add Azure DevOps provider commands to multitool CLI (projects, pipelines, repos, boards)
        2. Implement Azure DevOps REST API integration for real operations
        3. Add simulation mode support via cube-server proxy for DevOps resources
        4. Integrate with existing Azure functions for unified Azure management
        5. Add comprehensive authentication and organization handling for Azure DevOps
        6. Document Azure DevOps multitool usage and CI/CD integration examples
- [x] Add Azure functions to multitool framework (high priority) - COMPLETED:
    - [x] Implement Azure Monitor, Log Analytics, Application Insights, and Budget management
    - [x] Support for AKS cluster operations and monitoring stack creation
    - [x] Integration with existing transformer/generator Azure functions
    - [x] Step-by-step suggestions:
        1. [x] Add Azure provider commands to multitool CLI (monitor, loganalytics, appinsights, budget, aks)
        2. [ ] Implement Azure SDK integration for real cloud operations - IN PROGRESS
        3. [x] Add simulation mode support via cube-server proxy
        4. [x] Integrate with multitool provider commands for Azure monitoring and budget operations
        5. [x] Add comprehensive error handling and validation for Azure operations
        6. [ ] Document Azure multitool usage and workflow examples - PENDING
- [ ] Test Azure functions with simulator and direct mode (high priority):
    - [ ] Test simulation mode functionality via cube-server
    - [ ] Implement missing API client methods (ExecuteSimulation, GetProviderInfo, etc.)
    - [ ] Test direct mode with real Azure SDK integration
    - [ ] Validate error handling and output formatting
    - [ ] Add integration tests for both modes
- [ ] Enhance werfty-generator application (high priority):
    - Creates new Terraform templates from cloud state (JSON/YAML → .tf files)
    - Multi-cloud code generation (Azure, AWS, GCP, Hetzner, StackIT, IONOS)
    - Integration with cube-server for simulation workflows
    - Uses shared/ models for consistent data structures
    - Step-by-step suggestions:
        1. Extend werfty-generator to support more cloud providers and resources
        2. Add JSON/YAML input support for template generation from multitool exports
        3. Implement provider-specific template optimization and best practices
        4. Add CLI interface for batch template generation
        5. Integrate with multitool workflow (mt export → werfty-generator → .tf files)
        6. Use shared/ models for input/output data consistency
        7. Document werfty-generator usage and integration examples
- [ ] Enhance werfty-transformator application (high priority):
    - Imports and exports Terraform code between cloud providers (.tf ↔ .tf)
    - Cross-cloud migration and translation capabilities
    - Uses shared/ models for provider mapping and translation rules
    - Integration with cube-server for validation and testing
    - Step-by-step suggestions:
        1. Implement HCL parsing and AST manipulation for .tf files
        2. Add provider mapping rules using shared/ models
        3. Support bidirectional transformation (Azure ↔ AWS, GCP ↔ Hetzner, etc.)
        4. Add validation via cube-server simulation before transformation
        5. Implement resource dependency analysis and preservation
        6. Add CLI interface for batch transformations
        7. Document transformation rules and supported provider pairs
- [ ] Develop terraform-multicloud-provider (high priority):
    - Generic Terraform provider for simulation and multi-cloud workflows
    - Proxies operations to cube-server for simulation
    - Supports real cloud operations via provider delegation
    - Uses shared/ models for consistent resource definitions
    - Step-by-step suggestions:
        1. [x] Scaffold terraform-multicloud-provider basic structure
        2. [ ] Implement Terraform provider plugin architecture using shared/ models
        3. [ ] Add resource CRUD operations that proxy to cube-server for simulation
        4. [ ] Support real cloud operations via direct provider delegation
        5. [ ] Add migration interface for werfty-transformator HCL conversion
        6. [ ] Implement provider mapping and simulation workflows
        7. [ ] Add comprehensive integration tests for simulation and migration
        8. [ ] Document provider usage, simulation, and migration workflows
- [ ] Enhance multitool CLI capabilities (high priority):
    - Direct cloud operations with option to proxy via cube-server
    - Export cloud state for werfty-generator input
    - Import/export configurations for werfty-transformator
    - Unified object storage management across all clouds
    - Step-by-step suggestions:
        1. Add export commands to generate JSON/YAML for werfty-generator
        2. Add import commands to apply configurations from werfty-transformator
        3. Extend object storage commands (create, list, delete, etc.) for all providers
        4. Implement provider selection and config (direct vs. proxy/simulation)
        5. Add batch operations and scripting support
        6. Integrate with sim-server/cube-server APIs for simulation mode
        7. Add tests and documentation for all supported providers and modes

## Next Steps
- [ ] Complete shared/ module integration and standardization
- [ ] Fix compilation issues in client/pkg/output/formatter.go
- [ ] Test multitool Azure functions in both direct and simulation modes
- [ ] Implement missing API client methods for cube-server integration
- [ ] Add comprehensive integration tests for the complete workflow:
    1. multitool export cloud state → JSON/YAML
    2. werfty-generator create Terraform → .tf files
    3. werfty-transformator migrate between clouds → .tf files
    4. terraform-provider-punchbag simulate and validate
- [ ] Set up automated integration testing (CI, example validation)
- [ ] Enhance generator/transformator for more providers/resources and config-driven workflows
- [ ] Expand documentation and developer onboarding materials
- [ ] Fix and re-enable failing provider simulation API tests in `server/api/provider_simulation_test.go`
- [ ] Add automated Go test for Hetzner S3 bucket create/delete workflow using automation mode flags
- [ ] Document and automate the steps to spin up and tear down Hetzner Kubernetes (K8s) service

## Integration Workflow Vision
1. **multitool export**: Download current cloud state → JSON/YAML files
2. **werfty-generator**: Convert cloud state → Terraform .tf files
3. **werfty-transformator**: Migrate .tf files between cloud providers
4. **terraform-multicloud-provider**: Simulate and test Terraform configurations
5. **cube-server**: Provide simulation backend for all components
6. **shared/**: Ensure consistent data models across all applications

## Low Priority / Future
- [ ] Integration workflow between generator and backend/server
- [ ] Add more advanced provider transformation features as needed
- [ ] Add support for StackIT Object Storage in generator
- [ ] Further automate the Terraform generation workflow (e.g., config-driven, batch generation, etc.)
- [ ] Add more provider pairs and conversion logic to `werfty-transformator`
- [ ] Update all source files at the root level to include an SPDX license comment for AGPL-3.0-only

## Notes
- All new features and refactors should use the unified model layer in `shared/`
- Keep Terraform provider code and backend logic strictly separated
- Ensure clear distinction between:
    - **werfty-generator**: Creates new Terraform code from cloud state
    - **werfty-transformator**: Transforms existing Terraform code between providers
    - **multitool**: Direct cloud operations with optional cube-server proxy
    - **terraform-provider-punchbag**: Generic provider for simulation workflows
- Always check shared/ module usage for consistency across applications
- Use environment variables (e.g., PUNCHBAG_BASE_DIR) for all scripts and automation

## Architecture Validation Checklist
- [ ] All applications properly use shared/ models
- [ ] Clear separation between generator (new code) and transformator (existing code)
- [ ] multitool supports both direct and proxy modes
- [ ] cube-server integration works across all components
- [ ] API consistency between multitool, generators, and terraform provider
