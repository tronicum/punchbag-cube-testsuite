package azure

import (
	"fmt"
)

func CreateAzureResources() {
	fmt.Println("Creating "azurerm_log_analytics_workspace": "example"")
	fmt.Println("Creating "azurerm_application_insights": "example"")
	fmt.Println("Creating "azurerm_monitor_metric_alert": "example"")
	fmt.Println("Creating "azurerm_eventhub_namespace": "example"")
	fmt.Println("Creating "azurerm_monitor_diagnostic_setting": "example"")
}
