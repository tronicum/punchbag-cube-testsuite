package generator

import "fmt"

func GenerateAzureMonitoringTemplate(config map[string]interface{}) string {
	return "# Azure Monitoring Template\nresource \"azurerm_monitoring\" \"example\" {}"
}

func GenerateAzureKubernetesTemplate(config map[string]interface{}) string {
	return "# Azure Kubernetes Template\nresource \"azurerm_kubernetes_cluster\" \"example\" {}"
}

func GenerateAzureBudgetTemplate(config map[string]interface{}) string {
	return "# Azure Budget Template\nresource \"azurerm_budget\" \"example\" {}"
}
