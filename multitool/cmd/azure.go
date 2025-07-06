package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
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

// AKS management commands
var azureAksCmd = &cobra.Command{
	Use:   "aks",
	Short: "Manage Azure Kubernetes Service (AKS) clusters",
}

var azureAksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an AKS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		projectID, _ := cmd.Flags().GetString("project-id")
		api := client.NewAPIClient(proxyServer)
		clusterClient := client.NewClusterClient(api)
		req := &models.ClusterCreateRequest{
			Name:          name,
			Provider:      models.Azure,
			ResourceGroup: resourceGroup,
			Location:      location,
			ProjectID:     projectID,
		}
		cluster, err := clusterClient.CreateCluster(req)
		if err != nil {
			fmt.Printf("Failed to create AKS cluster: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("AKS cluster created: %s (ID: %s)\n", cluster.Name, cluster.ID)
	},
}

var azureAksGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get AKS cluster details",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		clusterClient := client.NewClusterClient(api)
		cluster, err := clusterClient.GetCluster(id)
		if err != nil {
			fmt.Printf("Failed to get AKS cluster: %v\n", err)
			os.Exit(1)
		}
		b, _ := json.MarshalIndent(cluster, "", "  ")
		fmt.Println(string(b))
	},
}

var azureAksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List AKS clusters",
	Run: func(cmd *cobra.Command, args []string) {
		api := client.NewAPIClient(proxyServer)
		clusterClient := client.NewClusterClient(api)
		clusters, err := clusterClient.ListClustersByProvider(models.Azure)
		if err != nil {
			fmt.Printf("Failed to list AKS clusters: %v\n", err)
			os.Exit(1)
		}
		b, _ := json.MarshalIndent(clusters, "", "  ")
		fmt.Println(string(b))
	},
}

var azureAksDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an AKS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		clusterClient := client.NewClusterClient(api)
		err := clusterClient.DeleteCluster(id)
		if err != nil {
			fmt.Printf("Failed to delete AKS cluster: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("AKS cluster deleted: %s\n", id)
	},
}

// Log Analytics management commands
var azureLogCmd = &cobra.Command{
	Use:   "loganalytics",
	Short: "Manage Azure Log Analytics workspaces",
}

var azureLogCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Log Analytics workspace",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating Log Analytics workspace (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureLogGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Log Analytics workspace details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting Log Analytics workspace (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureLogListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Log Analytics workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing Log Analytics workspaces (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureLogDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Log Analytics workspace",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting Log Analytics workspace (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

// Azure Budget management commands
var azureBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage Azure Budgets",
}

var azureBudgetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Azure Budget",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating Azure Budget (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureBudgetGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Azure Budget details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting Azure Budget (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureBudgetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Azure Budgets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing Azure Budgets (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureBudgetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Azure Budget",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting Azure Budget (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

// Application Insights management commands
var azureAppInsightsCmd = &cobra.Command{
	Use:   "appinsights",
	Short: "Manage Azure Application Insights instances",
}

var azureAppInsightsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Application Insights instance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating Application Insights instance (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureAppInsightsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Application Insights instance details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting Application Insights instance (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureAppInsightsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Application Insights instances",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing Application Insights instances (not yet implemented)")
		// TODO: Implement real API/proxy call
	},
}

var azureAppInsightsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Application Insights instance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting Application Insights instance (not yet implemented)")
		// TODO: Implement real API/proxy call
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

	azureAksCmd.AddCommand(azureAksCreateCmd, azureAksGetCmd, azureAksListCmd, azureAksDeleteCmd)
	azureCmd.AddCommand(azureAksCmd)

	azureLogCmd.AddCommand(azureLogCreateCmd, azureLogGetCmd, azureLogListCmd, azureLogDeleteCmd)
	azureCmd.AddCommand(azureLogCmd)

	azureBudgetCmd.AddCommand(azureBudgetCreateCmd, azureBudgetGetCmd, azureBudgetListCmd, azureBudgetDeleteCmd)
	azureCmd.AddCommand(azureBudgetCmd)

	azureCmd.AddCommand(azureGetMonitorCmd)
	azureCmd.AddCommand(azureGetLogAnalyticsCmd)

	azureAppInsightsCmd.AddCommand(azureAppInsightsCreateCmd, azureAppInsightsGetCmd, azureAppInsightsListCmd, azureAppInsightsDeleteCmd)
	azureCmd.AddCommand(azureAppInsightsCmd)
}
