# k8sctl Documentation

This document covers usage, commands, and architecture for the k8sctl CLI subcommand.

## Overview
- `k8sctl` is a top-level subcommand of the multitool CLI for kubectl-like operations.
- All operations are provider-agnostic and work through the multitool abstraction.

## Usage Examples
```sh
# Apply a manifest
./multitool/mt k8sctl apply -f manifest.yaml --profile aws-dev

# Get cluster info
./multitool/mt k8sctl get cluster --provider hetzner --profile hetzner-prod
```

## Developer Notes
- All provider logic must use the shared/ library abstraction.
- No direct provider/model code outside shared/.
- UX mirrors kubectl, but always operates through multitool.

## References
- See main README.md for high-level project info and links to other app docs.
- See `shared/README.md` for shared library API and usage.
- See `multitool/README.md` for CLI documentation.

---
SPDX-License-Identifier: AGPL-3.0-only
