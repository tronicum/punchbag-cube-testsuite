# Punchbag Cube Test Suite

A comprehensive multi-cloud test suite for testing punchbag cube functionality with server, werfty, and Terraform provider components.

## Overview

This project provides a complete ecosystem for testing various aspects of the punchbag cube system across multiple cloud providers including:
- REST API server for multi-cloud cluster management and test execution
- Command-line werfty for interacting with the API across cloud providers
- Terraform provider for Infrastructure as Code (IaC) support
- Multi-cloud support: Azure (AKS), StackIT (Schwarz IT), AWS (EKS), GCP (GKE)
- Performance and load testing capabilities across different cloud environments

## Supported Cloud Providers

- **Azure**: Azure Kubernetes Service (AKS)
- **StackIT (Schwarz IT)**: StackIT Kubernetes Engine (SKE)
- **AWS**: Amazon Elastic Kubernetes Service (EKS) [planned]
- **GCP**: Google Kubernetes Engine (GKE) [planned]

## Project Structure

```
├── README.md                    # This file
├── LICENSE                      # License file
├── server/                      # REST API Server (Multi-Cloud)
│   ├── main.go                  # Server entry point
│   ├── go.mod                   # Server dependencies
│   ├── Dockerfile               # Server container config
│   ├── README.md                # Server documentation
│   ├── api/                     # API layer
│   │   ├── handlers.go          # Multi-cloud HTTP request handlers
│   │   ├── routes.go            # Route definitions
│   │   └── openapi.yaml         # API specification
│   ├── models/                  # Data models
│   │   └── aks.go               # Multi-cloud cluster models
│   └── store/                   # Data storage layer
│       └── store.go             # Multi-cloud storage interface
├── werfty/                      # CLI Werfty (Multi-Cloud)
│   ├── main.go                  # Werfty entry point
│   ├── go.mod                   # Werfty dependencies
│   ├── README.md                # Werfty documentation
│   └── pkg/                     # Werfty packages
│       ├── api/                 # API werfty
│       │   └── werfty.go        # Multi-cloud HTTP werfty
│       └── output/              # Output formatting
│           └── formatter.go     # Multi-cloud output formatters
├── terraform-provider/          # Terraform Provider (Multi-Cloud)
│   ├── main.go                  # Provider entry point
│   ├── go.mod                   # Provider dependencies (includes StackIT provider)
│   ├── README.md                # Provider documentation
│   └── internal/                # Provider implementation
│       └── provider/            # Multi-cloud provider logic
├── models/                      # Shared data models
│   └── aks.go                   # Multi-cloud cluster models
└── store/                       # Shared storage layer
    └── store.go                 # Multi-cloud storage interface
```
            ├── provider.go      # Main provider
            ├── cluster_resource.go
            ├── test_resource.go
            └── clusters_data_source.go
```

## Components

### 🖥️ Server (`/server`)

REST API server that provides endpoints for:
- AKS cluster management (CRUD operations)
- Test execution and monitoring
- Health checks and metrics
- Complete OpenAPI specification

**Quick Start:**
```bash
cd server
go run main.go
```

Server will be available at `http://localhost:8080`

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
cd server
docker build -t punchbag-server .
docker run -p 8080:8080 punchbag-server

# Werfty (for CI/CD pipelines)
cd werfty
docker build -t punchbag-werfty .
```

## API Documentation

The server provides comprehensive API documentation:
- Interactive docs: `http://localhost:8080/docs`
- OpenAPI spec: `server/api/openapi.yaml`

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

## 🏃 Local Development & Testing

### Start the API Server
```bash
# From project root
cd cmd/cube-server
# Or use Go workspace mode
PORT=8081 go run main.go
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

## Contributing

Each component has its own documentation and development guidelines. See the README.md files in each directory for specific instructions.

## License

This project is licensed under the terms specified in the LICENSE file.
