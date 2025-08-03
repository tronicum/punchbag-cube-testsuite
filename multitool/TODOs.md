# 2025-08-03: CLI Help and Command Consistency

- [ ] Update multitool CLI help output to accurately reflect all available commands (e.g., `local os-detect`).
- [ ] Ensure all documented examples in README and help match the actual CLI structure and subcommands.
aws-s3# TODO: Support generic AWS S3 simulation endpoints
- Add CLI support for --storage-provider generic-aws-s3 to test against the generic AWS S3 simulation endpoint
- Add tests for generic AWS S3 simulation mode
## Modularization and Migration (from root TODOs.md)

- [ ] Move Hetzner S3 simulation mock (`NewHetznerS3Mock` and related code) from `multitool/pkg/client/hetzner_s3_mock.go` to `shared/providers/hetzner/objectstorage.go`.
- [ ] Move the CLI command `simulate-hetzner-s3` from `multitool/cmd/sim_hetzner_s3.go` to `cube-server/cmd/sim_hetzner_s3.go`.
- [ ] Register the simulation command in cube-server only, and update documentation to reflect the new location.
- [ ] Remove legacy simulation code from multitool after migration.

## Notes on Hetzner S3 Bucket Metadata
- [ ] Hetzner Object Storage (and all S3-compatible APIs) do not provide bucket creation or update timestamps via the S3 API. This is a limitation of the protocol and not the implementation. If richer metadata is needed, monitor Hetzner's hcloud API for future support.
# Multitool TODOs and Documentation

This file contains all multitool-specific tasks, configuration notes, usage examples, and migration notes. For overall framework information, see the root TODOs.md.

## Configuration and CLI Flag Precedence

The multitool CLI supports flexible configuration for all commands, including `k8sctl` and `k8s-manage`. The following precedence is used for flags such as `--mode` and `--provider`:

1. **CLI flag** (e.g. `--mode`, `--provider`)
2. **Environment variable** (`K8SCTL_MODE`, `K8SCTL_PROVIDER`)
3. **User config** (`$HOME/.mt/config.yaml`)
4. **Project config** (`./conf/k8sctl.yml`)
5. **Default** (hardcoded fallback)

This allows you to set global, per-user, or per-project defaults, and override them at runtime as needed.

**Example for k8sctl:**

```sh
# Use local mode by default (set in conf/k8sctl.yml or $HOME/.mt/config.yaml)
mt k8sctl get nodes

# Override mode for a single command
mt k8sctl get nodes --mode=proxy

# Use an environment variable for a session
export K8SCTL_MODE=direct
mt k8sctl get pods
```

### --provider flag meaning

The `--provider` flag is context-sensitive:

- For `k8s-manage` and other cloud lifecycle commands, it refers to the cloud provider (e.g., `hetzner`, `azure`, `aws`, `gcp`, etc.).
- For `k8sctl`, it may refer to the Kubernetes provider context, which can be mapped to a specific kubeconfig or cluster abstraction.

Always check the command help for the expected values and usage.

#### Example config files

**conf/k8sctl.yml**
```yaml
default_mode: local
default_provider: hetzner
```

**$HOME/.mt/config.yaml**
```yaml
default_mode: proxy
default_provider: azure
```

## Migration Notes

- The multitool CLI documentation has been moved to `multitool/README.md`.

## Multitool CLI and Configuration

See [multitool/README.md](./multitool/README.md) for full documentation on the multitool CLI, configuration system, usage examples, and migration notes.

## Multitool-Specific TODOs

- [x] Modularize simulation logic: Hetzner S3 mock migrated to shared/providers/hetzner. Legacy code removed from multitool. Simulation endpoints exposed via cube-server REST API.
- [x] Refactor object storage logic to use shared models and interfaces (`shared/models`, `shared/providers/hetzner/objectstorage.go`). All CLI and server code now use these abstractions for provider operations.
- [x] Implement per-profile config loading (`multitool/.mtconfig/<profile>/config.yaml`). CLI flag and env var support for profile selection is available.
- [x] Add/expand tests for all simulation endpoints, including error cases and edge conditions. CI checks for module hygiene and shared usage are in place and passing.
- [ ] Use shared/log and shared/errors for all error handling and logging in multitool and sim server.
- [ ] Improve help output, usage examples, and shell completion scripts. Add batch and interactive modes for resource operations.
- [ ] Update README and TODOs to reflect new architecture, usage, and migration notes.
- [ ] Add support for injecting dummy S3 buckets via ENV (e.g., SIMULATE_DUMMY_S3_BUCKETS) in the simulation server for test setup. Must support all providers: aws-s3, generic-aws-s3, hetzner.
- [ ] Document and enforce: all test orchestration must be done via multitool CLI, never by direct file manipulation or bash logic. Provider names must be used exactly as: aws-s3, generic-aws-s3, hetzner.

---
SPDX-License-Identifier: AGPL-3.0-only
