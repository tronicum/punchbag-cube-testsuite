package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/client"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// azureCmd is the parent for Azure-specific commands
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Manage Azure resources",
}

// azureGetMonitorCmd fetches Azure Monitor resource state
var azureGetMonitorCmd = &cobra.Command{
	Use:   "get monitor",
	Short: "Download Azure Monitor resource state as JSON",
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		name, _ := cmd.Flags().GetString("name")
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = fmt.Sprintf("monitor_%s.json", name)
		}
		var url string
		if proxyServer != "" {
			url = fmt.Sprintf("%s/api/v1/azure/monitor?resource_group=%s&name=%s", proxyServer, resourceGroup, name)
		} else {
			url = fmt.Sprintf("https://management.azure.com/subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/%s/providers/Microsoft.Insights/monitors/%s?api-version=2021-09-01", resourceGroup, name)
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Failed to fetch Azure Monitor: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		f, err := os.Create(output)
		if err != nil {
			fmt.Printf("Failed to write file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		enc.Encode(data)
		fmt.Printf("Azure Monitor state saved to %s\n", output)
	},
}

// azureGetLogAnalyticsCmd fetches Azure Log Analytics resource state
var azureGetLogAnalyticsCmd = &cobra.Command{
	Use:   "get log-analytics",
	Short: "Download Azure Log Analytics resource state as JSON",
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		name, _ := cmd.Flags().GetString("name")
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = fmt.Sprintf("loganalytics_%s.json", name)
		}
		var url string
		if proxyServer != "" {
			url = fmt.Sprintf("%s/api/v1/azure/loganalytics?resource_group=%s&name=%s", proxyServer, resourceGroup, name)
		} else {
			url = fmt.Sprintf("https://management.azure.com/subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s?api-version=2021-12-01-preview", resourceGroup, name)
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Failed to fetch Azure Log Analytics: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		f, err := os.Create(output)
		if err != nil {
			fmt.Printf("Failed to write file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		enc.Encode(data)
		fmt.Printf("Azure Log Analytics state saved to %s\n", output)
	},
}

// --- Log Analytics Management Commands ---
var azureCreateLogAnalyticsCmd = &cobra.Command{
	Use:   "create log-analytics",
	Short: "Create an Azure Log Analytics workspace",
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		name, _ := cmd.Flags().GetString("name")
		location, _ := cmd.Flags().GetString("location")
		retention, _ := cmd.Flags().GetInt("retention-days")
		apiClient := client.NewAPIClient(proxyServer)
		logClient := client.NewLogAnalyticsClient(apiClient)
		workspace := &sharedmodels.LogAnalyticsWorkspace{
			Name:          name,
			ResourceGroup: resourceGroup,
			Location:      location,
			RetentionDays: retention, // Fixed field name
		}
		result, err := logClient.Create(workspace)
		if err != nil {
			fmt.Printf("Failed to create Log Analytics workspace: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created Log Analytics workspace: %s (ID: %s)\n", result.Name, result.ID)
	},
}

var azureListLogAnalyticsCmd = &cobra.Command{
	Use:   "list log-analytics",
	Short: "List Azure Log Analytics workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		apiClient := client.NewAPIClient(proxyServer)
		logClient := client.NewLogAnalyticsClient(apiClient)
		workspaces, err := logClient.List()
		if err != nil {
			fmt.Printf("Failed to list Log Analytics workspaces: %v\n", err)
			os.Exit(1)
		}
		for _, ws := range workspaces {
			fmt.Printf("- %s (ID: %s, Group: %s, Location: %s)\n", ws.Name, ws.ID, ws.ResourceGroup, ws.Location)
		}
	},
}

var azureDeleteLogAnalyticsCmd = &cobra.Command{
	Use:   "delete log-analytics",
	Short: "Delete an Azure Log Analytics workspace",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		apiClient := client.NewAPIClient(proxyServer)
		logClient := client.NewLogAnalyticsClient(apiClient)
		err := logClient.Delete(id)
		if err != nil {
			fmt.Printf("Failed to delete Log Analytics workspace: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted Log Analytics workspace: %s\n", id)
	},
}

// --- Application Insights Management Commands ---
var azureCreateAppInsightsCmd = &cobra.Command{
	Use:   "create appinsights",
	Short: "Create an Azure Application Insights resource",
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		name, _ := cmd.Flags().GetString("name")
		location, _ := cmd.Flags().GetString("location")
		appType, _ := cmd.Flags().GetString("app-type")
		retention, _ := cmd.Flags().GetInt("retention-days")
		apiClient := client.NewAPIClient(proxyServer)
		appClient := client.NewAppInsightsClient(apiClient)
		app := &sharedmodels.AppInsightsResource{
			Name:          name,
			ResourceGroup: resourceGroup,
			Location:      location,
			AppType:       appType,   // Use appType field
			RetentionDays: retention, // Use retention field
			// InstrumentationKey and other fields can be set if needed
		}
		result, err := appClient.Create(app)
		if err != nil {
			fmt.Printf("Failed to create App Insights resource: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created App Insights resource: %s (ID: %s)\n", result.Name, result.ID)
	},
}

var azureListAppInsightsCmd = &cobra.Command{
	Use:   "list appinsights",
	Short: "List Azure Application Insights resources",
	Run: func(cmd *cobra.Command, args []string) {
		apiClient := client.NewAPIClient(proxyServer)
		appClient := client.NewAppInsightsClient(apiClient)
		apps, err := appClient.List()
		if err != nil {
			fmt.Printf("Failed to list App Insights resources: %v\n", err)
			os.Exit(1)
		}
		for _, app := range apps {
			fmt.Printf("- %s (ID: %s, Group: %s, Location: %s)\n", app.Name, app.ID, app.ResourceGroup, app.Location)
		}
	},
}

var azureDeleteAppInsightsCmd = &cobra.Command{
	Use:   "delete appinsights",
	Short: "Delete an Azure Application Insights resource",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		apiClient := client.NewAPIClient(proxyServer)
		appClient := client.NewAppInsightsClient(apiClient)
		err := appClient.Delete(id)
		if err != nil {
			fmt.Printf("Failed to delete App Insights resource: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted App Insights resource: %s\n", id)
	},
}

func init() {
	azureGetMonitorCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureGetMonitorCmd.Flags().String("name", "", "Azure Monitor resource name")
	azureGetMonitorCmd.Flags().String("output", "", "Output file (default: monitor_<name>.json)")
	azureGetMonitorCmd.MarkFlagRequired("resource-group")
	azureGetMonitorCmd.MarkFlagRequired("name")

	azureGetLogAnalyticsCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureGetLogAnalyticsCmd.Flags().String("name", "", "Azure Log Analytics workspace name")
	azureGetLogAnalyticsCmd.Flags().String("output", "", "Output file (default: loganalytics_<name>.json)")
	azureGetLogAnalyticsCmd.MarkFlagRequired("resource-group")
	azureGetLogAnalyticsCmd.MarkFlagRequired("name")

	azureCreateLogAnalyticsCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateLogAnalyticsCmd.Flags().String("name", "", "Log Analytics workspace name")
	azureCreateLogAnalyticsCmd.Flags().String("location", "", "Azure location")
	azureCreateLogAnalyticsCmd.Flags().Int("retention-days", 30, "Retention days (default: 30)")
	azureCreateLogAnalyticsCmd.MarkFlagRequired("resource-group")
	azureCreateLogAnalyticsCmd.MarkFlagRequired("name")
	azureCreateLogAnalyticsCmd.MarkFlagRequired("location")

	azureDeleteLogAnalyticsCmd.Flags().String("id", "", "Log Analytics workspace ID")
	azureDeleteLogAnalyticsCmd.MarkFlagRequired("id")

	azureCreateAppInsightsCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateAppInsightsCmd.Flags().String("name", "", "App Insights resource name")
	azureCreateAppInsightsCmd.Flags().String("location", "", "Azure location")
	azureCreateAppInsightsCmd.Flags().String("app-type", "web", "App type (default: web)")
	azureCreateAppInsightsCmd.Flags().Int("retention-days", 90, "Retention days (default: 90)")
	azureCreateAppInsightsCmd.MarkFlagRequired("resource-group")
	azureCreateAppInsightsCmd.MarkFlagRequired("name")
	azureCreateAppInsightsCmd.MarkFlagRequired("location")

	azureDeleteAppInsightsCmd.Flags().String("id", "", "App Insights resource ID")
	azureDeleteAppInsightsCmd.MarkFlagRequired("id")

	azureCmd.AddCommand(azureGetMonitorCmd)
	azureCmd.AddCommand(azureGetLogAnalyticsCmd)
	azureCmd.AddCommand(azureCreateLogAnalyticsCmd, azureListLogAnalyticsCmd, azureDeleteLogAnalyticsCmd)
	azureCmd.AddCommand(azureCreateAppInsightsCmd, azureListAppInsightsCmd, azureDeleteAppInsightsCmd)
}
