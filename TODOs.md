# TODOs for punchbag-cube-testsuite

## Project Architecture Overview

**Core Components (All use shared/ library):**
- **shared/**: Central library containing all cloud operations, models, and utilities
  - Azure, AWS, GCP, Hetzner, StackIT, IONOS providers
  - Common data models and interfaces
  - Cube-server integration and simulation logic
  - Cloud state export/import functionality
- **multitool/mt**: CLI for direct cloud operations (uses shared/ library)
- **werfty-generator**: CLI for creating Terraform from cloud state (uses shared/ library)
- **werfty-transformator**: CLI for transforming Terraform between providers (uses shared/ library)
- **terraform-multicloud-provider**: Terraform provider (uses shared/ library)
- **cube-server**: Simulation backend (uses shared/ library)
- **client/**: API client utilities (part of shared ecosystem)

## Current Sprint Tasks (This Week)

### Phase A: Shared Library Integration (CRITICAL)
- [ ] Move Azure commands from multitool to shared/ library
- [ ] Create shared/providers/azure/ package with all Azure operations
- [ ] Update multitool to use shared/providers/azure instead of local commands
- [ ] Standardize shared/ library API for all cloud providers
- [ ] Add shared/export and shared/import packages for data exchange
- [ ] Ensure all applications use only shared/ library for cloud operations

### Phase B: Strategic Architecture (AFTER A)
- [x] Scaffold terraform-multicloud-provider structure
- [ ] Update terraform-provider to use shared/ library
- [ ] Design provider plugin architecture using shared/ models

### Phase C: Test & Validate (AFTER A & B)
- [ ] Test all applications using shared/ library
- [ ] Validate loose coupling between components
- [ ] Test export/import workflows via shared/ library

## Milestones

### IMMEDIATE: Shared Library Migration (CRITICAL PRIORITY)
- [ ] Migrate Azure cloud operations to shared/ library:
    - Move Azure Monitor, Log Analytics, Application Insights, Budget operations to shared/providers/azure/
    - Move AKS cluster operations to shared/providers/azure/
    - Create shared/providers/azure/monitor.go, shared/providers/azure/aks.go, etc.
    - Update multitool to import and use shared/providers/azure
- [ ] Create shared/export and shared/import packages:
    - Export cloud state to JSON/YAML (used by werfty-generator)
    - Import configurations and validate (used by werfty-transformator)
    - Common data exchange formats across all applications
- [ ] Standardize shared/ library interfaces:
    - Common Provider interface for all clouds
    - Standardized authentication and configuration
    - Unified error handling and logging
    - Consistent simulation vs direct mode handling
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
- [x] Add Azure functions to multitool framework (COMPLETED) - **NEEDS MIGRATION TO SHARED/**:
    - [x] Implement Azure operations (MOVE TO shared/providers/azure/)
    - [x] Support for AKS cluster operations (MOVE TO shared/providers/azure/)
    - [ ] **CRITICAL: Migrate all Azure code from multitool to shared/ library**
    - Step-by-step migration:
        1. Create shared/providers/azure/ package structure
        2. Move Azure commands from multitool/cmd/azure.go to shared/providers/azure/
        3. Update multitool to import and use shared/providers/azure
        4. Create shared/export package for cloud state export
        5. Create shared/import package for configuration import
        6. Test multitool using shared/ library instead of local code

- [ ] Enhance werfty-generator application (high priority):
    - **Uses shared/export for cloud state input**
    - **Uses shared/providers/* for cloud operations**
    - Creates new Terraform templates from cloud state (JSON/YAML → .tf files)
    - Step-by-step suggestions:
        1. Use shared/export to read cloud state data
        2. Use shared/providers/azure for Azure-specific template generation
        3. Use shared/models for consistent data structures
        4. Add provider-agnostic template generation using shared/ interfaces
        5. Support all cloud providers via shared/providers/*
        6. Document werfty-generator usage with shared/ library

- [ ] Enhance werfty-transformator application (high priority):
    - **Uses shared/import for configuration validation**
    - **Uses shared/providers/* for cloud validation**
    - Transforms Terraform code between cloud providers (.tf ↔ .tf)
    - Step-by-step suggestions:
        1. Use shared/import to validate configurations
        2. Use shared/providers/* for cloud-specific transformations
        3. Use shared/models for provider mapping and translation rules
        4. Implement HCL parsing with shared/ validation
        5. Support bidirectional transformation using shared/ interfaces
        6. Document transformation workflows with shared/ library

- [ ] Enhance multitool CLI capabilities (high priority):
    - **Uses shared/providers/* for all cloud operations**
    - **Uses shared/export and shared/import for data exchange**
    - Independent CLI tool - no direct dependencies on other applications
    - Step-by-step suggestions:
        1. **Migrate to use shared/providers/* instead of local commands**
        2. Use shared/export to generate data for other applications
        3. Use shared/import to consume configurations
        4. Implement all cloud providers via shared/providers/*
        5. Add batch operations using shared/ library
        6. Test independence from other applications

## Next Steps (Immediate Actions Required)

### 1. Create Shared Library Structure
```
shared/
├── providers/
│   ├── azure/
│   │   ├── monitor.go      # Azure Monitor operations
│   │   ├── aks.go          # AKS operations
│   │   ├── budget.go       # Budget operations
│   │   └── client.go       # Azure SDK integration
│   ├── aws/
│   ├── gcp/
│   └── interface.go        # Common Provider interface
├── export/
│   ├── json.go            # JSON export functionality
│   └── yaml.go            # YAML export functionality
├── import/
│   ├── config.go          # Configuration import
│   └── validate.go        # Validation functionality
├── models/
│   ├── cluster.go         # Common cluster models
│   └── resource.go        # Common resource models
└── simulation/
    └── client.go          # Cube-server integration
```

### 2. Migration Priority Order
1. **Create shared/providers/azure/** and migrate Azure code from multitool
2. **Update multitool** to use shared/providers/azure
3. **Create shared/export** for werfty-generator integration
4. **Create shared/import** for werfty-transformator integration
5. **Test all applications** using shared/ library independently

### 3. Validation Checklist
- [ ] multitool works independently using shared/ library
- [ ] werfty-generator works independently using shared/ library  
- [ ] werfty-transformator works independently using shared/ library
- [ ] No direct dependencies between applications (only via shared/)
- [ ] All cloud operations centralized in shared/providers/*
- [ ] Data exchange happens via shared/export and shared/import

## Integration Workflow Vision (Via Shared Library)
1. **multitool**: Uses shared/providers/* → shared/export → JSON/YAML files
2. **werfty-generator**: Uses shared/export to read JSON/YAML → shared/providers/* → .tf files
3. **werfty-transformator**: Uses shared/import to read .tf → shared/providers/* → .tf files
4. **terraform-multicloud-provider**: Uses shared/providers/* and shared/simulation
5. **cube-server**: Uses shared/simulation for all components
6. **All data flows through shared/ library - no direct component dependencies**

## Low Priority / Future
- [ ] Integration workflow between generator and backend/server
- [ ] Add more advanced provider transformation features as needed
- [ ] Add support for StackIT Object Storage in generator
- [ ] Further automate the Terraform generation workflow (e.g., config-driven, batch generation, etc.)
- [ ] Add more provider pairs and conversion logic to `werfty-transformator`
- [ ] Update all source files at the root level to include an SPDX license comment for AGPL-3.0-only

## Notes
- **multitool is the central hub for all cloud operations** - generator and transformer depend on it
- All Azure commands live in multitool and are consumed by other components
- Generator uses multitool export commands to get cloud state data
- Transformer uses multitool validation commands to test transformations
- All new features and refactors should use the unified model layer in `shared/`
- Keep Terraform provider code and backend logic strictly separated
- Ensure clear distinction between:
    - **multitool**: Central cloud operations CLI (used by all other components)
    - **werfty-generator**: Creates new Terraform code (consumes multitool exports)
    - **werfty-transformator**: Transforms existing Terraform code (uses multitool validation)
    - **terraform-multicloud-provider**: Generic provider for simulation workflows
- Always check shared/ module usage for consistency across applications
- Use environment variables (e.g., PUNCHBAG_BASE_DIR) for all scripts and automation

## Architecture Validation Checklist
- [ ] **ALL Azure code moved to shared/providers/azure/**
- [ ] **multitool uses shared/ library exclusively**
- [ ] **werfty-generator uses shared/ library exclusively**
- [ ] **werfty-transformator uses shared/ library exclusively**
- [ ] No direct dependencies between applications
- [ ] All cloud operations in shared/providers/*
- [ ] Data exchange via shared/export and shared/import
- [ ] Each application can be used independently as CLI tool
- [ ] Consistent interfaces across all shared/ packages
