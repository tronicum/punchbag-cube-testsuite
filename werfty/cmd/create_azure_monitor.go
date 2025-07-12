package cmd

import (
	"fmt"
	"os"
	"time"

	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createAzureMonitorCmd = &cobra.Command{
	Use:   "create-azure-monitor",
	Short: "Create Azure Monitor services for Kubernetes clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))

		// Extended Azure Monitor services to include all logging, analytics, and related products
		services := []string{
			"log-analytics",
			"application-insights",
			"metrics",
			"event-hubs",
			"diagnostic-settings",
			"monitor-alerts",
			"container-insights",
			"service-map",
			"vm-insights",
			"network-watcher",
		}
		responses := make(map[string]interface{})

		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		workspaceName, _ := cmd.Flags().GetString("workspace-name")

		for _, service := range services {
			response := map[string]interface{}{
				"service":            service,
				"status":             "success",
				"resource_id":        fmt.Sprintf("azure-monitor-%s-%d", service, time.Now().Unix()),
				"resource_group":     resourceGroup,
				"location":           location,
				"workspace_name":     workspaceName,
				"timestamp":          time.Now().Format(time.RFC3339),
				"terraform_template": generateAzureMonitorTerraform(service, resourceGroup, location, workspaceName),
			}
			responses[service] = response
		}

		output.PrintJSON(responses, os.Stdout)
	},
}

func generateAzureMonitorTerraform(service, resourceGroup, location, workspaceName string) string {
	switch service {
	case "log-analytics":
		return fmt.Sprintf(`
resource "azurerm_log_analytics_workspace" "%s" {
  name                = "%s"
  location            = "%s"
  resource_group_name = "%s"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}`, workspaceName, workspaceName, location, resourceGroup)
	case "application-insights":
		return fmt.Sprintf(`
resource "azurerm_application_insights" "%s-ai" {
  name                = "%s-ai"
  location            = "%s"
  resource_group_name = "%s"
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.%s.id
}`, workspaceName, workspaceName, location, resourceGroup, workspaceName)
	case "container-insights":
		return fmt.Sprintf(`
resource "azurerm_log_analytics_solution" "container-insights" {
  solution_name         = "ContainerInsights"
  location              = "%s"
  resource_group_name   = "%s"
  workspace_resource_id = azurerm_log_analytics_workspace.%s.id
  workspace_name        = azurerm_log_analytics_workspace.%s.name
  
  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/ContainerInsights"
  }
}`, location, resourceGroup, workspaceName, workspaceName)
	default:
		return fmt.Sprintf("# Azure Monitor %s template placeholder", service)
	}
}

func init() {
	rootCmd.AddCommand(createAzureMonitorCmd)

	createAzureMonitorCmd.Flags().String("resource-group", "", "Azure resource group")
	createAzureMonitorCmd.Flags().String("location", "eastus", "Azure location")
	createAzureMonitorCmd.Flags().String("workspace-name", "", "Log Analytics workspace name")
	createAzureMonitorCmd.MarkFlagRequired("resource-group")
	createAzureMonitorCmd.MarkFlagRequired("workspace-name")
}
