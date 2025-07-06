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
resource "azurerm_monitor_metric_alert" "example" {
  name                = "test-monitor"
  resource_group_name = "example-rg"
  severity            = 3
  criteria            = ""
  // ...map more fields from JSON as needed
}