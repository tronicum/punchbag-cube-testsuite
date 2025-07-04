# Dummy Terraform template for Azure services

provider "azurerm" {
  features {}
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-log-analytics"
  location            = "East US"
  resource_group_name = "example-resources"
  sku                 = "PerGB2018"
}

resource "azurerm_application_insights" "example" {
  name                = "example-app-insights"
  location            = "East US"
  resource_group_name = "example-resources"
  application_type    = "web"
}

resource "azurerm_monitor_metric_alert" "example" {
  name                = "example-metric-alert"
  resource_group_name = "example-resources"
  scopes              = [azurerm_log_analytics_workspace.example.id]
  criteria {
    metric_namespace = "Microsoft.Insights"
    metric_name      = "requests"
    aggregation      = "Total"
    operator         = "GreaterThan"
    threshold        = 100
  }
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-eventhub"
  location            = "East US"
  resource_group_name = "example-resources"
  sku                 = "Standard"
}

resource "azurerm_monitor_diagnostic_setting" "example" {
  name               = "example-diagnostic-setting"
  target_resource_id = azurerm_log_analytics_workspace.example.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
  logs {
    category = "Administrative"
    enabled  = true
    retention_policy {
      enabled = true
      days    = 30
    }
  }
}
