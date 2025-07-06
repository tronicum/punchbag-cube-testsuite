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

### 📦 Package Management (Coming Soon)
- Cross-platform package installation and management
- OS detection and appropriate package manager selection
- Support for Homebrew, apt, rpm, pacman, choco, and winget

## Dual-Mode Operation

Multitool can operate in two modes:

- **Direct mode (default):** Calls real cloud provider APIs directly for resource management.
- **Proxy mode (with `--server` flag):** Forwards all resource management requests to a cube-server instance, which can simulate or proxy to real cloud providers.

### Example Usage

**Direct (real cloud):**
```sh
multitool cluster create my-cluster azure --resource-group my-rg --location eastus
```

**Proxy (via cube-server):**
```sh
multitool --server http://localhost:8080 cluster create my-cluster azure --resource-group my-rg --location eastus
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
multitool config init
```

### 2. Connect to Server (for Real Operations)

```bash
# Set server URL
multitool config set server_url http://your-server:8080

# Create a real cluster (requires server)
multitool cluster create prod-cluster azure --resource-group prod-rg --location eastus

# List clusters
multitool cluster list

# Run a test
multitool test run cluster-id connectivity
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

## 🚀 Terraform Deployment Automation

Multitool now supports unified Terraform workflows for all clouds:

```sh
multitool deploy plan <tf-file>
multitool deploy apply
multitool deploy destroy <tf-file>
```

- Use Werfty or compatible generator to create Terraform files for Azure, AWS, or GCP.
- Then use `multitool deploy` to plan, apply, or destroy resources.

### Example
```sh
go run generator/main.go generate -i examples/example_azure_storage_account.yaml -o storage_account.tf -p azure
multitool deploy plan storage_account.tf
multitool deploy apply
```

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
