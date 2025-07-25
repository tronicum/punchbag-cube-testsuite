# Multitool - Multi-Cloud CLI for Resource Management

Multitool is a comprehensive command-line interface for managing cloud resources, running tests, and handling system operations across multiple cloud providers. It supports Azure, AWS, GCP, Hetzner, IONOS, and StackIT.

## Features

### ☀️ Multi-Cloud Cluster Management
- Create, list, get, and delete Kubernetes clusters across multiple cloud providers
- Support for Azure AKS, AWS EKS, Google GKE, and other managed Kubernetes services
- Provider-specific configuration and resource management

### 🧪 Test Management
- Run various types of tests on clusters (connectivity, performance, security, compliance)
- Track test results and metrics
- Integration with the punchbag server for centralized test management

### ⚙️ Configuration Management
- Profile-based configuration for different environments
- Support for multiple output formats (table, JSON, YAML)
- Centralized credential and setting management

### ☁️ Object Storage Management
- Create, list, get, and delete S3-like object storage buckets across all supported cloud providers (AWS, Azure, GCP, StackIT, Hetzner, IONOS)
- Works in both direct (mock/local) and proxy (simulation/server) modes
- Provider is selected by the first argument to each objectstorage command (e.g., `mt objectstorage create my-bucket us-west-2 --provider aws`)
- In proxy mode (`--server`), requests are forwarded to the cube-server API for simulation or real provider proxying

## Configuration System

Multitool supports a kubeconfig-like profile system for cloud provider configuration.

- Profiles are stored in `.mtconfig/<profile>/config.yaml`.
- Switch profiles using `--profile <name>` or set the `MTCONFIG_PROFILE` environment variable.
- Each config file supports provider, region, credentials, endpoints, and custom settings.
- Environment variables in config values (e.g. `${AWS_ACCESS_KEY_ID}`) are expanded at runtime. If not set, multitool will attempt to read from a file named after the variable, or prompt the user.

### Example config.yaml
```yaml
provider: aws
region: eu-central-1
credentials:
  access_key: ${AWS_ACCESS_KEY_ID}
  secret_key: ${AWS_SECRET_ACCESS_KEY}
endpoints:
  s3: https://s3.eu-central-1.amazonaws.com
settings:
  versioning: true
  custom_policy: ./policy.json
```

## CLI Flag and Config Precedence

1. CLI flag (e.g. `--provider`, `--profile`)
2. Environment variable (e.g. `MTCONFIG_PROFILE`)
3. User config (`$HOME/.mt/config.yaml`)
4. Project config (`./conf/k8sctl.yml`)
5. Default (hardcoded fallback)

## Example Usage

```sh
# Use a specific profile
mt objectstorage list --profile aws-dev

# Use environment variable for profile
export MTCONFIG_PROFILE=default
mt objectstorage create mybucket eu-central-1
```

## Migration Notes

- The new config system replaces previous ad-hoc config files and flags.
- See the main README.md for a summary and links to provider-specific docs.

## Building the CLI

To build the CLI as `mt`:

```sh
cd multitool
# Build the binary as 'mt'
go build -o mt
```

You can now use `./mt` for all commands, e.g.:

```sh
./mt cluster list
./mt objectstorage create aws my-bucket us-west-2
```

## Dual-Mode Operation

mt can operate in two modes:

- **Direct mode (default):** Calls real cloud provider APIs directly for resource management.
- **Proxy mode (with `--server` flag):** Forwards all resource management requests to a cube-server instance, which can simulate or proxy to real cloud providers.

### Example Usage

**Direct (real cloud):**
```sh
mt cluster create my-cluster azure --resource-group my-rg --location eastus
```

**Proxy (via cube-server):**
```sh
mt --server http://localhost:8080 cluster create my-cluster azure --resource-group my-rg --location eastus
```
- If the cube-server is in simulation mode, this simulates the operation.
- If the cube-server is in real mode, this proxies the request to the real cloud provider.

## Installation

### Build from Source

```bash
git clone <repository-url>
cd multitool
go build -o multitool
```

## Note

Simulation and mock storage are only available in the cube-server and werfty/generator. Multitool is for real cloud management only.

## Quick Start

### 1. Initialize Configuration

```bash
mt config init
```

### 2. Connect to Server (for Real Operations)

```bash
# Set server URL
mt config set server_url http://your-server:8080

# Create a real cluster (requires server)
mt cluster create prod-cluster azure --resource-group prod-rg --location eastus

# List clusters
mt cluster list

# Run a test
mt test run cluster-id connectivity
```

## Command Reference

### Cluster Management

#### Create a Cluster

```bash
# Azure
multitool cluster create my-cluster azure \
  --resource-group my-rg \
  --location eastus \
  --config cluster-config.json

# AWS
multitool cluster create my-cluster aws \
  --region us-west-2 \
  --config cluster-config.json

# GCP
multitool cluster create my-cluster gcp \
  --project-id my-project \
  --region us-central1 \
  --config cluster-config.json
```

#### List Clusters

```bash
# List all clusters
multitool cluster list

# List clusters by provider
multitool cluster list azure

# Different output formats
multitool cluster list --output json
multitool cluster list --output yaml
multitool cluster list --output table
```

#### Get Cluster Details

```bash
multitool cluster get cluster-id
```

#### Delete a Cluster

```bash
multitool cluster delete cluster-id
multitool cluster delete cluster-id --confirm  # Skip confirmation
```

### Test Management

#### Run Tests

```bash
# Run connectivity test
multitool test run cluster-id connectivity

# Run performance test with config
multitool test run cluster-id performance --config perf-config.json

# Run security test
multitool test run cluster-id security

# Run compliance test
multitool test run cluster-id compliance
```

#### List Test Results

```bash
# List all test results for a cluster
multitool test list cluster-id

# Get specific test result
multitool test get test-id
```

### Configuration Management

#### Initialize Configuration

```bash
multitool config init
multitool config init --force  # Overwrite existing config
```

#### View Configuration

```bash
# Show all configuration
multitool config show

# Show specific profile
multitool config show development
multitool config show production
```

#### Set Configuration Values

```bash
# Global settings
multitool config set server_url http://localhost:8080
multitool config set default_provider azure
multitool config set default_region eastus
multitool config set default_output table

# Profile-specific settings
multitool config set --profile development provider aws
multitool config set --profile development region us-west-2
multitool config set --profile production server_url https://prod-server.com
```

#### List Profiles

```bash
multitool config list-profiles
```

### Legacy Commands

#### Kubernetes Commands (Deprecated)

```bash
# Use new cluster commands instead
multitool k8s get azure
multitool k8s create azure
multitool k8s delete azure
```

#### Package Management

```bash
# Detect OS and package manager
multitool os-detect

# List installed packages
multitool list-packages

# Install package (shows command, doesn't execute)
multitool install-package docker
```

## Object Storage Management

Multitool supports unified object storage management for all supported clouds. The provider is always specified as the first argument to each objectstorage command:

```sh
# Create a bucket in AWS (direct/mock mode)
multitool objectstorage create aws my-bucket us-west-2

# Create a bucket in Azure (direct/mock mode)
multitool objectstorage create azure my-bucket westeurope

# List buckets in GCP (direct/mock mode)
multitool objectstorage list gcp

# Use proxy/simulation mode (forward to cube-server)
multitool --server http://localhost:8080 objectstorage create stackit my-bucket eu01
multitool --server http://localhost:8080 objectstorage list hetzner
```

- The provider argument must be one of: `aws`, `azure`, `gcp`, `stackit`, `hetzner`, `ionos`.
- In direct/mock mode, multitool uses local mock logic for all operations.
- In proxy mode (with `--server`), multitool forwards all object storage commands to the specified server, which can simulate or proxy to real providers.
- You can set a default provider in your config, but the provider argument is always required for objectstorage commands.

## Configuration File

The configuration file is stored at `~/.multitool/config.yaml` and supports:

```yaml
server_url: http://localhost:8080
default_provider: azure
default_region: eastus
default_output: table

profiles:
  default:
    server_url: http://localhost:8080
    provider: azure
    region: eastus
    
  development:
    server_url: http://localhost:8080
    provider: azure
    region: eastus
    resource_group: dev-rg
    tags:
      environment: development
      team: devops
      
  production:
    server_url: https://punchbag-prod.example.com
    provider: azure
    region: eastus
    resource_group: prod-rg
    tags:
      environment: production
      team: devops
```

## Output Formats

Multitool supports three output formats:

### Table (Default)
```bash
multitool cluster list --output table
```
```
ID               Name             Provider  Status   Region       Created
---              ----             --------  ------   ------       -------
cluster-123      my-cluster       azure     running  eastus       2025-01-01 10:00
cluster-456      test-cluster     aws       running  us-west-2    2025-01-01 11:00
```

### JSON
```bash
multitool cluster list --output json
```
```json
[
  {
    "id": "cluster-123",
    "name": "my-cluster",
    "provider": "azure",
    "status": "running",
    "region": "eastus",
    "created_at": "2025-01-01T10:00:00Z"
  }
]
```

### YAML
```bash
multitool cluster list --output yaml
```
```yaml
- id: cluster-123
  name: my-cluster
  provider: azure
  status: running
  region: eastus
  created_at: 2025-01-01T10:00:00Z
```

## Supported Cloud Providers

- **Azure**: Azure Kubernetes Service (AKS)
- **AWS**: Amazon Elastic Kubernetes Service (EKS)
- **GCP**: Google Kubernetes Engine (GKE)
- **Hetzner**: Hetzner Cloud Kubernetes
- **IONOS**: IONOS Cloud Kubernetes
- **StackIT**: StackIT Kubernetes Engine

## Test Types

- **connectivity**: Test cluster connectivity and network access
- **performance**: Load testing and performance metrics
- **security**: Security scanning and vulnerability assessment
- **compliance**: Compliance checks and policy validation

## Integration with Punchbag Server

Multitool is designed to work with the punchbag server for centralized resource management:

1. **Server API**: All cluster and test operations can be performed through the server API
2. **Centralized Storage**: Cluster metadata and test results are stored on the server
3. **Multi-user Support**: Multiple users can share clusters and test results
4. **Audit Trail**: All operations are logged on the server

### Server Requirements

To use real cluster operations (non-simulation), you need:

1. A running punchbag server
2. Network connectivity to the server
3. Appropriate credentials configured

## Development

### Project Structure

```
multitool/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and CLI setup
│   ├── cluster.go         # Cluster management commands
│   ├── config.go          # Configuration management
│   └── ...
├── pkg/                   # Shared packages
│   ├── client/            # API client for server communication
│   ├── models/            # Data models (shared with server)
│   └── output/            # Output formatting
├── main.go               # CLI entry point
├── go.mod                # Go module definition
└── README.md             # This file
```

### Adding New Commands

1. Create a new file in `cmd/` directory
2. Implement the command using Cobra
3. Add the command to `root.go`
4. Update this README

### Adding New Cloud Providers

1. Add the provider to `models/cluster.go`
2. Update validation in cluster commands
3. Add provider-specific configuration handling
4. Update documentation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit a pull request

## License

[Add your license information here]

## Support

For issues and questions:

1. Check the documentation
2. Search existing issues
3. Create a new issue with detailed information
4. Include command output and error messages

## Roadmap

- [ ] Package management implementation
- [ ] More cloud provider support
- [ ] Kubernetes resource management (beyond clusters)
- [ ] Enhanced test types and metrics
- [ ] Configuration templates and presets
- [ ] Shell completion scripts
- [ ] Integration with CI/CD pipelines

# Modularization & Release Workflow

This project uses Go workspaces to modularize the CLI (`mt`) and shared Go modules. To build, test, and release:

- Use `go.work` at the repo root for local development across modules.
- Each module (e.g., `multitool/`, `shared/`) has its own `go.mod`.
- All import paths use canonical module paths (e.g., `github.com/tronicum/punchbag-cube-testsuite/shared/models`).

## Building & Testing

```sh
# Build the CLI
cd multitool
go build -o mt

# Run all tests in multitool and shared modules
cd multitool && go test ./...
cd ../shared && go test ./...
```

## Release

- To release the CLI, build the binary as above and publish it to your preferred distribution channel.
- To release shared Go modules, tag the commit and push to GitHub. Consumers can use `go get` with the tag.

# Azure Logging & Application Insights Management

The CLI supports managing Azure Log Analytics and Application Insights resources:

- Create, list, and delete Log Analytics workspaces:
  - `mt azure create log-analytics ...`
  - `mt azure list log-analytics`
  - `mt azure delete log-analytics --id <id>`
- Create, list, and delete Application Insights resources:
  - `mt azure create appinsights ...`
  - `mt azure list appinsights`
  - `mt azure delete appinsights --id <id>`

These commands are structured for easy extension to other providers in the future.

# CI/CD

- Set up GitHub Actions or another CI/CD system to automate building, testing, and releasing the CLI and modules.
- Example workflow: build and test on push/PR, release on tag.
