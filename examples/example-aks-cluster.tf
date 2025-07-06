# Example Terraform configuration for creating an AKS cluster

provider "azurerm" {
  features {}
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks-cluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  default_node_pool {
    name       = "default"
    node_count = 3
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  kubernetes_version = "1.24.0"
  dns_prefix         = "exampleaks"

  tags = {
    Environment = "example"
  }
}

output "kube_config" {
  value = azurerm_kubernetes_cluster.example.kube_config_raw
}

# Dummy Azure credentials for local testing and CI
# These are NOT valid for real deployments. Replace with real values for production.
#
# You can override these by setting environment variables or a terraform.tfvars file.
#
# Example:
#   export TF_VAR_azure_subscription_id=your-subscription-id
#   export TF_VAR_azure_client_id=your-client-id
#   export TF_VAR_azure_client_secret=your-client-secret
#   export TF_VAR_azure_tenant_id=your-tenant-id

# See variables.tf for dummy defaults.
