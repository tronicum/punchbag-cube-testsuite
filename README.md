# Configuration and CLI Flag Precedence

The multitool CLI supports flexible configuration for all commands, including `k8sctl` and `k8s-manage`. The following precedence is used for flags such as `--mode` and `--provider`:

1. **CLI flag** (e.g. `--mode`, `--provider`)
2. **Environment variable** (`K8SCTL_MODE`, `K8SCTL_PROVIDER`)
3. **User config** (`$HOME/.mt/config.yaml`)
4. **Project config** (`./conf/k8sctl.yml`)
5. **Default** (hardcoded fallback)

This allows you to set global, per-user, or per-project defaults, and override them at runtime as needed.

**Example for k8sctl:**

## Migration Notes

- The multitool CLI documentation has been moved to `multitool/README.md`.
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

## Multitool CLI and Configuration

See [multitool/README.md](./multitool/README.md) for full documentation on the multitool CLI, configuration system, usage examples, and migration notes.

## CLI Structure

A comprehensive multi-cloud test suite for testing punchbag cube functionality with server, werfty, and Terraform provider components.

### ⚙️ Configuration & Object Storage Management
See [multitool/README.md](./multitool/README.md) for full details on configuration profiles, object storage, and advanced CLI usage.
## Supported Cloud Providers


- **AWS, GCP, etc.**: Extendable via shared/ abstractions and simulation handlers

Example usage:

```bash
curl -X POST http://localhost:8081/api/simulate/azure/loganalytics

```


- **Azure**: Azure Kubernetes Service (AKS), Azure Monitor, Log Analytics, Application Insights, Azure Budgets
- **StackIT (Schwarz IT)**: StackIT Kubernetes Engine (SKE)
- **AWS**: Amazon Elastic Kubernetes Service (EKS), CloudWatch, AWS Budgets [planned]
- **GCP**: Google Kubernetes Engine (GKE), Google Cloud Monitoring [planned]
- **Hetzner Cloud**: Hetzner Kubernetes, Hetzner Cloud Monitoring
- **IONOS Cloud**: IONOS Kubernetes, IONOS Monitoring

## Azure Features

### 🔵 Azure Kubernetes Service (AKS)
- Create, manage, and monitor AKS clusters
- Integrated Azure Monitor for containers
- Cost management with Azure Budgets
- Terraform generation for AKS resources

### 📊 Azure Monitor Integration
- Log Analytics workspace creation and management
- Application Insights for application monitoring
- Container Insights for Kubernetes monitoring
- Service Map and VM Insights
- Network Watcher integration

### 💰 Azure Cost Management
- Azure Budget creation with alerting
- Cost monitoring and reporting
- Integration with AKS cluster creation
- Terraform templates for budget management

## Azure Usage Examples

### Creating AKS Cluster with Monitoring and Budget

```bash
# Create AKS cluster with full monitoring and budget
./punchbag-werfty cluster create \
  --name my-aks-cluster \
  --provider azure \
  --resource-group my-rg \
  --location eastus \
  --kubernetes-version 1.28.0 \
  --node-count 3 \
  --enable-monitoring \
  --enable-budget \
  --budget-amount 2000.0

# Create Azure Monitor stack
./punchbag-werfty provider azure monitor create \
  --resource-group my-rg \
  --location eastus \
  --workspace-name my-monitoring-workspace

# Create Azure Budget
./punchbag-werfty provider azure budget create \
  --name my-project-budget \
  --amount 1500.0 \
  --resource-group my-rg \
  --time-grain Monthly \
  --alert-threshold 80.0
```

### Using Multitool for Azure Operations

```bash
# Create Azure monitoring stack
mt azure create monitoring-stack \
  --resource-group my-rg \
  --name my-monitoring \
  --location eastus

# Create Azure budget with monitoring
mt azure create budget-stack \
  --resource-group my-rg \
  --name my-budget \
  --amount 1000.0

# Download Azure Monitor configuration
mt azure get monitor \
  --resource-group my-rg \
  --name my-monitor \
  --output monitor_config.json

# Create Log Analytics workspace
mt azure create log-analytics \
  --resource-group my-rg \
  --name my-workspace \
  --location eastus \
  --retention-days 30
```

### Terraform Provider Examples

**Azure AKS cluster with monitoring:**
```hcl
resource "punchbag_cluster" "example_azure" {
  name     = "my-aks-cluster"
  provider = "azure"

  azure_config = {
    resource_group       = "my-rg"
    location            = "eastus"
    kubernetes_version  = "1.28.0"
    node_count          = 3
    enable_monitoring   = true
    enable_budget       = true
    budget_amount       = 2000.0
  }
}
```

**Azure Monitor stack:**
```hcl
resource "azurerm_log_analytics_workspace" "example" {
  name                = "my-workspace"
  location            = "eastus"
  resource_group_name = "my-rg"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "example" {
  name                = "my-ai"
  location            = "eastus"
  resource_group_name = "my-rg"
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.example.id
}
```

## 🔄 Werfty-Transformator: Enhanced Azure Support

The `werfty-transformator` now includes comprehensive Azure to AWS conversion:

### Azure to AWS Conversions
- Azure Monitor → AWS CloudWatch
- Azure Budget → AWS Budget
- Azure Log Analytics → CloudWatch Log Groups
- Application Insights → CloudWatch Dashboards
- Azure Blob Storage → AWS S3

### Usage Examples

```bash
# Convert Azure monitoring to AWS CloudWatch
go run main.go \
  --input azure_monitoring.tf \
  --src-provider azure \
  --destination-provider aws

# Convert complete Azure stack to AWS
go run main.go \
  --input azure_full_stack.tf \
  --src-provider azure \
  --destination-provider aws \
  --terraspace
```

## 🧪 API Usage Examples

### Azure-specific API calls:

```bash
# Create Azure Monitor services
curl -X POST http://localhost:8081/api/v1/azure/monitor \
  -H 'Content-Type: application/json' \
  -d '{
    "resource_group": "my-rg",
    "location": "eastus",
    "workspace_name": "my-workspace",
    "services": ["log-analytics", "application-insights", "container-insights"]
  }'

# Create Azure Budget
curl -X POST http://localhost:8081/api/v1/azure/budget \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "my-budget",
    "amount": 1000.0,
    "resource_group": "my-rg",
    "time_grain": "Monthly",
    "alert_threshold": 80.0
  }'
```

## Components

### 🖥️ Server (`cube-server/`)

Unified REST API server and simulation backend. All endpoints, simulation logic, and orchestration are now in `cube-server/`.

**Quick Start:**
```bash
cd cube-server
go mod tidy
go run main.go
```

Server will be available at `http://localhost:8081`

### 📱 Werfty (`/werfty`)

Command-line interface for interacting with the server:
- Manage clusters from the command line
- Run tests and monitor progress
- Multiple output formats (table, JSON, YAML)
- Configuration file support

**Quick Start:**
```bash
cd werfty

go build -o punchbag-werfty .

./punchbag-werfty cluster list
```

### 🏗️ Terraform Provider (`/terraform-provider`)

Terraform provider for Infrastructure as Code:
- Manage clusters with Terraform
- Run tests as part of infrastructure deployment
- Import existing resources
- Complete resource and data source support

**Quick Start:**
```bash
cd terraform-provider
go build -o terraform-provider-punchbag
```

## Getting Started

### Prerequisites
- Go 1.21 or later
- Docker (optional, for containerized deployment)
- Terraform 1.0+ (for using the provider)

### Quick Setup

1. **Start the Server:**
```bash
cd server
go mod tidy
go run main.go
```

2. **Use the Werfty:**
```bash
cd werfty
go mod tidy
go build -o punchbag-werfty .
./punchbag-werfty cluster list
```

3. **Try the Terraform Provider:**
```bash
cd terraform-provider
go mod tidy
go build -o terraform-provider-punchbag
# See terraform-provider/README.md for setup instructions
```

## Docker Support

Each component can be containerized:

```bash
# Server
cd cube-server
docker build -t punchbag-cube-server .
docker run -p 8081:8081 punchbag-cube-server

# Werfty (for CI/CD pipelines)
cd werfty
docker build -t punchbag-werfty .
```

## API Documentation

The server provides comprehensive API documentation:
- Interactive docs: `http://localhost:8081/docs`
- OpenAPI spec: `cube-server/api/openapi.yaml`

## Multi-Cloud Usage Examples

### Creating Clusters

**Azure (AKS) Cluster:**
```bash
./punchbag-werfty cluster create \
  --name my-aks-cluster \
  --provider azure \
  --resource-group my-rg \
  --location eastus \
  --kubernetes-version 1.28.0 \
  --node-count 3
```

**StackIT Cluster:**
```bash
./punchbag-werfty cluster create \
  --name my-stackit-cluster \
  --provider schwarz-stackit \
  --project-id your-project-id \
  --region eu-de-1 \
  --kubernetes-version 1.28.0 \
  --node-count 3
```

### Listing Clusters

**All clusters:**
```bash
./punchbag-werfty cluster list
```

**Filter by provider:**
```bash
./punchbag-werfty cluster list --provider azure
./punchbag-werfty cluster list --provider schwarz-stackit
```

### Running Tests

**Run tests on any cluster:**
```bash
./punchbag-werfty cluster test cluster-id-123 --type performance_test
```

### Terraform Provider Examples

**Azure cluster:**
```hcl
resource "punchbag_cluster" "example_azure" {
  name     = "my-aks-cluster"
  provider = "azure"

  azure_config = {
    resource_group     = "my-rg"
    location          = "eastus"
    kubernetes_version = "1.28.0"
    node_count        = 3
  }
}
```

**StackIT cluster:**
```hcl
resource "punchbag_cluster" "example_stackit" {
  name     = "my-stackit-cluster"
  provider = "schwarz-stackit"

  stackit_config = {
    project_id         = "your-project-id"
    region            = "eu-de-1"
    kubernetes_version = "1.28.0"
    node_count        = 3
  }
}
```

## 🔄 Werfty-Transformator: Terraform Conversion Tool

The `werfty-transformator` tool allows you to convert Terraform files between cloud providers (AWS, Azure, GCP) or to the generic multipass-cloud-layer provider.

### Usage

```bash
cd werfty-transformator

go run main.go --input <input.tf> --src-provider <azure|aws|gcp> --destination-provider <azure|aws|gcp|multipass-cloud-layer>
```

### Supported Conversions
- Azure Blob Storage ↔ AWS S3
- Any S3-like resource → multipass-cloud-layer

### Examples

**Convert Azure Blob Storage to AWS S3:**
```bash
go run main.go --input azure_blob_example.tf --src-provider azure --destination-provider aws
```

**Convert AWS S3 to Azure Blob Storage:**
```bash
go run main.go --input aws_s3_example.tf --src-provider aws --destination-provider azure
```

**Convert AWS S3 to multipass-cloud-layer:**
```bash
go run main.go --input aws_s3_example.tf --src-provider aws --destination-provider multipass-cloud-layer
```

### How it Works
- The tool parses the input Terraform file, detects supported resources, and rewrites them for the target provider.
- Only S3-like resources are supported for now, but the tool is extensible for more resource types and providers.

### Adding More Mappings
- Extend the `ConvertTerraform` function in `werfty-transformator/main.go` to add new provider/resource conversions.
- See the code comments for guidance.

## 🏃 Local Development & Testing

### Start the API Server
```bash
# From project root
cd cube-server
go mod tidy
go run main.go
```

### Run Automated API Tests
```bash
bash scripts/test_api.sh
```

### Run Linting
```bash
bash scripts/lint.sh
```

## 🧪 API Usage Examples

- Health check: `curl http://localhost:8081/health`
- List clusters: `curl http://localhost:8081/api/v1/clusters`
- Create cluster:
  ```bash
  curl -X POST http://localhost:8081/api/v1/clusters \
    -H 'Content-Type: application/json' \
    -d '{"name":"test-cluster","provider":"azure","location":"eastus"}'
  ```

## Testing GitHub Actions Workflows Locally

You can test your GitHub Actions workflows locally using [`act`](https://github.com/nektos/act):

1. Install act (macOS):
   ```sh
   brew install act
   ```
2. Run your workflow locally:
   ```sh
   act
   ```
   - Use `act pull_request` to simulate a pull request event.
   - Use `-j <job-name>` to run a specific job.
   - You may need Docker running and to set up secrets or environment variables for some jobs.

This is useful for quickly validating CI changes before pushing to GitHub.

## Contributing

Each component has its own documentation and development guidelines. See the README.md files in each directory for specific instructions.

## License

This project is licensed under the terms specified in the LICENSE file.
