# Werfty (Terraform Generator)

Werfty is a modular code generator and test suite for Terraform resources, focused on Azure and multicloud (AWS, GCP) infrastructure. It supports YAML/JSON config input, robust CLI workflows, and extensibility for new providers/resources.

## Features

- **YAML/JSON config input** for resource definitions
- **Cobra CLI** with subcommands: `generate`, `validate`, `simulate`
- **Generate Terraform code** for AKS, EKS, GKE, S3, Monitor, Log Analytics, App Insights, and more
- **Schema validation** for resource properties
- **Terraform output validation** (`terraform validate`, `tflint`) in tests and CI
- **Extensible**: add new providers/resources easily

## Quick Start

1. **Generate Terraform from YAML/JSON:**
   ```sh
   go run main.go generate -i examples/example_azure_services.yaml -o output.tf -p azure
   # or for AWS/GCP:
   go run main.go generate -i examples/example_aws_eks.json -o output.tf -p aws
   ```

2. **Validate a config file:**
   ```sh
   go run main.go validate -i examples/example_azure_services.yaml -p azure
   ```

3. **Simulate (dry-run) Terraform output:**
   ```sh
   go run main.go simulate -o output.tf
   # Runs terraform validate and tflint on the output
   ```

## Example YAML Config

```yaml
resourceType: aks
properties:
  name: my-aks
  location: eastus
  resourceGroup: my-rg
  nodeCount: 3
```

## Example JSON Config

```json
{
  "resourceType": "eks",
  "properties": {
    "name": "my-eks",
    "region": "us-west-2",
    "nodeCount": 2
  }
}
```

## Architecture & Extensibility

- All resource generators and providers are modular (see `internal/generator/`)
- Add new resources by implementing a generator and updating the provider
- CLI is fully extensible via Cobra subcommands

## Advanced Usage

- See `examples/` for more YAML/JSON configs
- See `docs/` for architecture diagrams and plugin/extensibility notes (coming soon)

## CI/CD & Deployment

- CI runs on every push and PR: lint, test, terraform validate, tflint.
- Deployment runs on push to `main` (see `.github/workflows/deploy.yml`).
- Store cloud credentials as GitHub Actions secrets (e.g., `AZURE_CLIENT_ID`, `AWS_ACCESS_KEY_ID`, etc.).
- Customize `deploy.sh` for your deployment needs (cloud, artifact, etc.).

Example secret usage in workflow:
```yaml
      - name: Set up Azure credentials
        run: |
          echo "AZURE_CLIENT_ID=${{ secrets.AZURE_CLIENT_ID }}" >> $GITHUB_ENV
          echo "AZURE_CLIENT_SECRET=${{ secrets.AZURE_CLIENT_SECRET }}" >> $GITHUB_ENV
          echo "AZURE_TENANT_ID=${{ secrets.AZURE_TENANT_ID }}" >> $GITHUB_ENV
```

---

For more, see the [full documentation](docs/) or run `go run main.go --help`.
