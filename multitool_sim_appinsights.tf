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
resource "azurerm_application_insights" "example" {
  name                = "test-appinsights"
  location            = "eastus"
  resource_group_name = "example-resource-group"
  application_type    = "web"
  retention_in_days   = 90
  // ...map more fields from JSON as needed
}