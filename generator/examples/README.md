# Examples Library

This directory contains real-world YAML and JSON config examples for Azure, AWS, GCP, and multicloud scenarios.

## Azure AKS Example (YAML)
- `example_azure_aks.yaml`: Minimal AKS cluster definition for Azure.

## AWS EKS Example (JSON)
- `example_aws_eks.json`: Minimal EKS cluster definition for AWS.

## GCP GKE Example (JSON)
- `example_gcp_gke.json`: Minimal GKE cluster definition for GCP.

## Multicloud Example (YAML)
- `example_multicloud.yaml`: Multiple resources across providers.

Use these with the CLI:
```sh
go run main.go generate -i examples/example_azure_aks.yaml -o output.tf -p azure
go run main.go generate -i examples/example_aws_eks.json -o output.tf -p aws
go run main.go generate -i examples/example_gcp_gke.json -o output.tf -p gcp
```
