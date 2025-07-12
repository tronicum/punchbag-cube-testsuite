# End-to-End Test: Generator and Server with Terraform Example

This test and example demonstrate how to use the generator and unified simulation API server with Terraform for Azure AKS.

## Prerequisites
- Terraform installed
- Go installed (for generator and server)

## Dummy Azure Credentials
This example uses dummy Azure credentials for local testing and CI. These are NOT valid for real deployments.

- `azure_subscription_id`: 00000000-0000-0000-0000-000000000000
- `azure_client_id`: 11111111-1111-1111-1111-111111111111
- `azure_client_secret`: dummy-secret
- `azure_tenant_id`: 22222222-2222-2222-2222-222222222222

Override these with real values using environment variables or a `terraform.tfvars` file:

```
export TF_VAR_azure_subscription_id=your-subscription-id
export TF_VAR_azure_client_id=your-client-id
export TF_VAR_azure_client_secret=your-client-secret
export TF_VAR_azure_tenant_id=your-tenant-id
```

## Steps

1. **Start the unified simulation server:**
   ```sh
   cd cube-server
   go run main.go
   ```

2. **Generate Terraform from JSON using the generator:**
   ```sh
   cd generator
   go run main.go --generate-terraform --input ../generator/test_aks_expanded.json --output ../examples/generated_aks.tf
   ```

3. **Initialize and validate Terraform in the examples directory:**
   ```sh
   cd ../examples
   terraform init
   terraform validate
   ```

4. **Run Terraform plan (will fail with dummy credentials, but validates the workflow):**
   ```sh
   terraform plan
   ```

## Notes
- The example-aks-cluster.tf file is ready for local and CI validation.
- For real Azure deployments, provide valid credentials as described above.
- The generator and server can be extended to support more resource types and scenarios.
