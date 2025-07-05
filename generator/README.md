# Werfty (Terraform Generator)

Werfty is a code generator for Terraform resources, focused on Azure (AKS, Budgets, Monitoring, etc.). It can simulate resources using the cube-server or generate real Terraform code for deployment.

## Features

- **Simulate resources** via cube-server REST API (no cloud calls required)
- **Generate Terraform code** for AKS clusters, Budgets, and more
- **Import existing resources** into Terraform state
- **Supports both simulation and real execution workflows**

## Quick Start

1. **Simulate an AKS cluster:**
   ```sh
   go run main.go --simulate-import --name my-aks --resource-group my-rg --location eastus --node-count 3
   ```
   (Requires cube-server running at http://localhost:8080)

2. **Generate Terraform for a new AKS cluster:**
   ```sh
   go run main.go --generate-terraform --name my-aks --resource-group my-rg --location eastus --node-count 3
   ```

3. **Import an existing AKS cluster into Terraform state:**
   ```sh
   terraform import azurerm_kubernetes_cluster.my_aks \
     /subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/my-rg/providers/Microsoft.ContainerService/managedClusters/my-aks
   ```

## API Credentials

- Set as environment variables (recommended):
  - `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID`
- Or use a `.env` file or config YAML for local testing

## Example Workflow

1. Simulate with Werfty → Review output
2. Generate Terraform code → Apply with `terraform apply`
3. Import existing resources with `terraform import` if needed

## Development

- Extend resource support by editing `main.go` and template functions
- Integrates with `shared/simulation` for simulation logic

## License

MIT
