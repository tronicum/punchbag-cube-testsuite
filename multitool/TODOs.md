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

- [ ] Implement `.mtconfig/<profile>/config.yaml` profile system for multitool
- [ ] Add CLI flag `--profile` and env var support
- [ ] Update provider selection logic to use config values
- [ ] Document usage and migration in README.md
- [ ] Expand multitool test coverage for all simulation endpoints (Azure, AWS, GCP, validation, etc.), including edge cases and error handling
- [ ] Add CI checks to enforce module hygiene, run tests, and validate shared usage
- [ ] Standardize error handling and logging using shared/log and shared/errors
- [ ] Improve help output and usage examples for all multitool commands
- [ ] Add shell completion scripts for multitool CLI
- [ ] Support batch resource operations from manifest files
- [ ] Add interactive mode for multitool CLI

---
SPDX-License-Identifier: AGPL-3.0-only
