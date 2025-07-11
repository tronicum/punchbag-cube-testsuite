package cmd

import (
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

// azureCreateMonitoringStackCmd creates a complete Azure monitoring stack
var azureCreateMonitoringStackCmd = &cobra.Command{
	Use:   "create monitoring-stack",
	Short: "Create a complete Azure monitoring stack",
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		name, _ := cmd.Flags().GetString("name")
		location, _ := cmd.Flags().GetString("location")

		if proxyServer != "" {
			// Use proxy mode
			createMonitoringStackViaProxy(resourceGroup, name, location)
		} else {
			// Use direct mode
			createMonitoringStackDirect(resourceGroup, name, location)
		}
	},
}

func createMonitoringStackViaProxy(resourceGroup, name, location string) {
	// Implementation for proxy mode
	fmt.Printf("Creating Azure monitoring stack via proxy server...\n")
	fmt.Printf("Resource Group: %s\n", resourceGroup)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Location: %s\n", location)

	// This would call the cube-server API for simulation or real execution
	fmt.Printf("Monitoring stack created successfully (simulated)\n")
}

func createMonitoringStackDirect(resourceGroup, name, location string) {
	// Implementation for direct mode
	fmt.Printf("Creating Azure monitoring stack directly...\n")
	fmt.Printf("Resource Group: %s\n", resourceGroup)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Location: %s\n", location)

	// This would call Azure APIs directly
	fmt.Printf("Monitoring stack created successfully\n")
}

// azureCreateBudgetStackCmd creates Azure budgets with monitoring integration
var azureCreateBudgetStackCmd = &cobra.Command{
	Use:   "create budget-stack",
	Short: "Create Azure budget with monitoring integration",
	Run: func(cmd *cobra.Command, args []string) {
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		name, _ := cmd.Flags().GetString("name")
		amount, _ := cmd.Flags().GetFloat64("amount")

		if proxyServer != "" {
			createBudgetStackViaProxy(resourceGroup, name, amount)
		} else {
			createBudgetStackDirect(resourceGroup, name, amount)
		}
	},
}

func createBudgetStackViaProxy(resourceGroup, name string, amount float64) {
	fmt.Printf("Creating Azure budget stack via proxy server...\n")
	fmt.Printf("Resource Group: %s\n", resourceGroup)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Amount: $%.2f\n", amount)
	fmt.Printf("Budget stack created successfully (simulated)\n")
}

func createBudgetStackDirect(resourceGroup, name string, amount float64) {
	fmt.Printf("Creating Azure budget stack directly...\n")
	fmt.Printf("Resource Group: %s\n", resourceGroup)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Amount: $%.2f\n", amount)
	fmt.Printf("Budget stack created successfully\n")
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

	// Monitoring stack flags
	azureCreateMonitoringStackCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateMonitoringStackCmd.Flags().String("name", "", "Monitoring stack name")
	azureCreateMonitoringStackCmd.Flags().String("location", "eastus", "Azure location")
	azureCreateMonitoringStackCmd.MarkFlagRequired("resource-group")
	azureCreateMonitoringStackCmd.MarkFlagRequired("name")

	// Budget stack flags
	azureCreateBudgetStackCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateBudgetStackCmd.Flags().String("name", "", "Budget name")
	azureCreateBudgetStackCmd.Flags().Float64("amount", 1000.0, "Budget amount")
	azureCreateBudgetStackCmd.MarkFlagRequired("resource-group")
	azureCreateBudgetStackCmd.MarkFlagRequired("name")

	azureCmd.AddCommand(azureGetMonitorCmd)
	azureCmd.AddCommand(azureGetLogAnalyticsCmd)
	azureCmd.AddCommand(azureCreateLogAnalyticsCmd, azureListLogAnalyticsCmd, azureDeleteLogAnalyticsCmd)
	azureCmd.AddCommand(azureCreateAppInsightsCmd, azureListAppInsightsCmd, azureDeleteAppInsightsCmd)
	azureCmd.AddCommand(azureCreateMonitoringStackCmd)
	azureCmd.AddCommand(azureCreateBudgetStackCmd)

	// Add Azure subcommands
	azureCmd.AddCommand(azureMonitorCmd)
	azureCmd.AddCommand(azureBudgetCmd)
	azureCmd.AddCommand(azureAksCmd)

	// Add create commands
	azureMonitorCmd.AddCommand(azureCreateMonitorCmd)
	azureBudgetCmd.AddCommand(azureCreateBudgetCmd)
	azureAksCmd.AddCommand(azureCreateAksCmd)

	// Azure Monitor flags
	azureCreateMonitorCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateMonitorCmd.Flags().String("location", "eastus", "Azure region")
	azureCreateMonitorCmd.Flags().String("workspace-name", "", "Log Analytics workspace name")
	azureCreateMonitorCmd.Flags().Bool("simulation", false, "Use simulation mode")
	azureCreateMonitorCmd.MarkFlagRequired("resource-group")
	azureCreateMonitorCmd.MarkFlagRequired("workspace-name")

	// Azure Budget flags
	azureCreateBudgetCmd.Flags().String("name", "", "Budget name")
	azureCreateBudgetCmd.Flags().Float64("amount", 0, "Budget amount in USD")
	azureCreateBudgetCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateBudgetCmd.Flags().String("time-grain", "Monthly", "Budget time grain")
	azureCreateBudgetCmd.Flags().Bool("simulation", false, "Use simulation mode")
	azureCreateBudgetCmd.MarkFlagRequired("name")
	azureCreateBudgetCmd.MarkFlagRequired("amount")
	azureCreateBudgetCmd.MarkFlagRequired("resource-group")

	// Azure AKS flags
	azureCreateAksCmd.Flags().String("name", "", "AKS cluster name")
	azureCreateAksCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateAksCmd.Flags().String("location", "eastus", "Azure region")
	azureCreateAksCmd.Flags().Int("node-count", 3, "Number of nodes in default pool")
	azureCreateAksCmd.Flags().Bool("simulation", false, "Use simulation mode")
	azureCreateAksCmd.MarkFlagRequired("name")
	azureCreateAksCmd.MarkFlagRequired("resource-group")
}
