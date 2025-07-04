# Terraform Provider for Punchbag Cube Test Suite

A Terraform provider for managing AKS clusters and tests in the Punchbag Cube Test Suite.

## Overview

This Terraform provider allows you to manage AKS clusters and run tests using Infrastructure as Code (IaC) practices. It provides resources and data sources for interacting with the Punchbag Cube Test Suite API.

## Features

- **Cluster Management**: Create, read, update, and delete AKS clusters
- **Test Execution**: Run tests on clusters via Terraform
- **Data Sources**: Query existing clusters and test results
- **Import Support**: Import existing resources into Terraform state

## Installation

### From Source

```bash
cd terraform-provider

# Build the provider
go build -o terraform-provider-punchbag

# Install locally for testing
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/punchbag/punchbag/1.0.0/darwin_amd64/
cp terraform-provider-punchbag ~/.terraform.d/plugins/registry.terraform.io/punchbag/punchbag/1.0.0/darwin_amd64/
```

### Using the Provider

```hcl
terraform {
  required_providers {
    punchbag = {
      source  = "registry.terraform.io/punchbag/punchbag"
      version = "~> 1.0"
    }
  }
}

provider "punchbag" {
  host     = "http://localhost:8080"
  username = "admin"
  password = "password"
}
```

## Resources

### `punchbag_cluster`

Manages an AKS cluster in the test suite.

```hcl
resource "punchbag_cluster" "example" {
  name               = "my-test-cluster"
  resource_group     = "test-rg"
  location           = "eastus"
  kubernetes_version = "1.28.0"
  node_count         = 3

  tags = {
    Environment = "test"
    Purpose     = "load-testing"
  }
}
```

#### Arguments

- `name` (Required) - Name of the AKS cluster
- `resource_group` (Required) - Azure resource group
- `location` (Required) - Azure location
- `kubernetes_version` (Optional) - Kubernetes version (default: "1.28.0")
- `node_count` (Optional) - Number of nodes (default: 3)
- `tags` (Optional) - Map of tags

#### Attributes

- `id` - Cluster identifier
- `status` - Current cluster status

### `punchbag_test`

Runs a test on an AKS cluster.

```hcl
resource "punchbag_test" "load_test" {
  cluster_id = punchbag_cluster.example.id
  test_type  = "load_test"

  config = {
    duration    = "5m"
    concurrency = "10"
    target_url  = "http://example.com"
  }
}
```

#### Arguments

- `cluster_id` (Required) - ID of the cluster to test
- `test_type` (Required) - Type of test (load_test, performance_test, stress_test)
- `config` (Optional) - Test configuration map

#### Attributes

- `id` - Test identifier
- `status` - Current test status

## Data Sources

### `punchbag_clusters`

Retrieves information about all clusters.

```hcl
data "punchbag_clusters" "all" {}

output "cluster_names" {
  value = [for cluster in data.punchbag_clusters.all.clusters : cluster.name]
}
```

## Examples

### Complete Example

```hcl
terraform {
  required_providers {
    punchbag = {
      source  = "registry.terraform.io/punchbag/punchbag"
      version = "~> 1.0"
    }
  }
}

provider "punchbag" {
  host = "http://localhost:8080"
}

# Create a cluster
resource "punchbag_cluster" "test_cluster" {
  name               = "terraform-test-cluster"
  resource_group     = "terraform-test-rg"
  location           = "eastus"
  kubernetes_version = "1.28.0"
  node_count         = 3

  tags = {
    Environment = "test"
    ManagedBy   = "terraform"
  }
}

# Run a load test on the cluster
resource "punchbag_test" "load_test" {
  cluster_id = punchbag_cluster.test_cluster.id
  test_type  = "load_test"

  config = {
    duration     = "10m"
    concurrency  = "20"
    request_rate = "100"
    target_url   = "http://my-app.example.com"
    method       = "GET"
  }
}

# Data source to list all clusters
data "punchbag_clusters" "all" {}

# Output cluster information
output "cluster_id" {
  value = punchbag_cluster.test_cluster.id
}

output "cluster_status" {
  value = punchbag_cluster.test_cluster.status
}

output "test_id" {
  value = punchbag_test.load_test.id
}

output "test_status" {
  value = punchbag_test.load_test.status
}

output "all_clusters" {
  value = data.punchbag_clusters.all.clusters
}
```

## Development

### Prerequisites

- Go 1.21 or later
- Terraform 1.0 or later

### Building

```bash
go mod tidy
go build -o terraform-provider-punchbag
```

### Testing

```bash
go test ./...
```

### Local Development

1. Build the provider
2. Create a local plugin directory
3. Copy the provider binary
4. Create a test Terraform configuration
5. Run `terraform init` and `terraform apply`

## Environment Variables

- `PUNCHBAG_HOST` - Server URL
- `PUNCHBAG_USERNAME` - Username for authentication
- `PUNCHBAG_PASSWORD` - Password for authentication
