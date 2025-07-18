# SPDX-License-Identifier: AGPL-3.0-only

# Architecture Overview: punchbag-cube-testsuite

## Copyright
Copyright (C) 2023-2025 tronicum@user.github.com

---

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

---

## Note on Google Cloud Code Generation

[Magic Modules](https://googlecloudplatform.github.io/magic-modules/) is an open-source code generation framework developed by Google. It is used to automatically generate infrastructure-as-code (IaC) provider code for Google Cloud Platform (GCP) resources, including Terraform, Ansible, and other tools. Magic Modules takes a set of YAML-based resource definitions and produces provider code, documentation, and tests, ensuring consistency and reducing manual maintenance for large cloud APIs.

**In this project:**
- Magic Modules will be used as the code generator for Google-related resources in the werfty-generator and werfty-transformator components.
- This enables automated, up-to-date support for GCP resources and simplifies the process of adding or updating Google provider features in Terraform and other supported formats.

See: https://googlecloudplatform.github.io/magic-modules/
