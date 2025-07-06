# Deployment Automation

## CI/CD Integration
- All generated Terraform is validated with `terraform validate` and `tflint` in CI.
- See `.github/workflows/ci.yml` for details.

## Cloud-Specific Deployment Scripts
- You can add scripts in `scripts/` for Azure, AWS, and GCP deployment automation.
- Example: `scripts/deploy_azure.sh` to run `terraform init`, `plan`, and `apply` for Azure resources.

## Example Azure Deployment Script
```sh
#!/bin/bash
set -e
cd $(dirname "$0")/..
TFFILE=${1:-output.tf}
RESOURCE_GROUP=${2:-my-rg}

terraform init
terraform plan -out=tfplan "$TFFILE"
terraform apply tfplan
```

## Advanced Usage
- Integrate with GitHub Actions for automated deployment after PR merge.
- Use environment variables for secrets and cloud credentials.
