#!/bin/bash
# Azure Terraform deployment automation script
set -e
cd $(dirname "$0")/..
TFFILE=${1:-output.tf}
RESOURCE_GROUP=${2:-my-rg}

terraform init
terraform plan -out=tfplan "$TFFILE"
terraform apply tfplan
