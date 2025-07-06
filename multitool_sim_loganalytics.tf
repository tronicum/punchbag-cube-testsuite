terraform {
  required_version = ">= 1.0.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}
resource "azurerm_log_analytics_workspace" "example" {
  name                = "test-loganalytics"
  location            = "eastus"
  resource_group_name = "example-resource-group"
  sku                 = "PerGB2018"
  retention_in_days   = 30
  // ...map more fields from JSON as needed
}