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

### `multipass_cloud_layer_bucket`

A generic, cross-cloud object storage bucket resource. This abstraction allows you to manage S3-like storage (AWS S3, Azure Blob, GCP GCS) using a unified resource.

```hcl
resource "multipass_cloud_layer_bucket" "example" {
  name         = "my-bucket"
  region       = "us-west-2"
  provider     = "aws" # or "azure", "gcp"
  storage_class = "STANDARD"
  tier         = "Hot" # For Azure Blob, optional
}
```

- `name`: The bucket/container name.
- `region`: The region for the bucket.
- `provider`: The cloud provider (`aws`, `azure`, or `gcp`).
- `storage_class`: (Optional) Storage class (e.g., `STANDARD`, `COOL`).
- `tier`: (Optional) Tier for Azure Blob or similar.

This resource will be mapped to the correct cloud-specific resource at apply time.

#### Example: AWS S3

```hcl
resource "multipass_cloud_layer_bucket" "aws_bucket" {
  name     = "my-s3-bucket"
  region   = "us-west-2"
  provider = "aws"
}
```

#### Example: Azure Blob

```hcl
resource "multipass_cloud_layer_bucket" "azure_blob" {
  name     = "myblob"
  region   = "westeurope"
  provider = "azure"
  tier     = "Hot"
}
```

#### Example: GCP GCS

```hcl
resource "multipass_cloud_layer_bucket" "gcs_bucket" {
  name     = "my-gcs-bucket"
  region   = "us-central1"
  provider = "gcp"
}
```

#### Example: IONOS S3

```hcl
resource "multipass_cloud_layer_bucket" "ionos_bucket" {
  name     = "my-ionos-bucket"
  region   = "de-fra"
  provider = "ionos"
}
```

#### Example: StackIT S3

```hcl
resource "multipass_cloud_layer_bucket" "stackit_bucket" {
  name     = "my-stackit-bucket"
  region   = "eu01"
  provider = "stackit"
}
```

#### Example: Hetzner Cloud S3 (hcloud)

```hcl
resource "multipass_cloud_layer_bucket" "hcloud_bucket" {
  name     = "my-hcloud-bucket"
  region   = "fsn1"
  provider = "hcloud"
}
```

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

# multipass-cloud-layer Terraform Provider

A cross-cloud object storage abstraction for AWS S3, Azure Blob, GCP GCS, IONOS S3, StackIT S3, and Hetzner Cloud S3.

## Publishing to the Terraform Registry

1. Ensure your provider code and documentation are up to date.
2. The provider address is set to `registry.terraform.io/multipass/cloud-layer`.
3. Tag a release in your VCS (e.g., GitHub) following [Terraform Registry requirements](https://developer.hashicorp.com/terraform/registry/providers/publishing).
4. Push your code and tag to the public repository.
5. The registry will automatically detect and index your provider.

## Example Usage

```hcl
terraform {
  required_providers {
    multipass-cloud-layer = {
      source  = "multipass/cloud-layer"
      version = ">= 0.1.0"
    }
  }
}

provider "multipass-cloud-layer" {
  host = "http://localhost:8080" # Your API/proxy endpoint
}

resource "multipass_cloud_layer_bucket" "example" {
  name     = "my-bucket"
  region   = "us-west-2"
  provider = "aws" # or "azure", "gcp", "ionos", "stackit", "hcloud"
}
```

See the main README for more details and examples.
