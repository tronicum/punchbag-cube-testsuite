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

	azureCmd.AddCommand(azureGetMonitorCmd)
	azureCmd.AddCommand(azureGetLogAnalyticsCmd)
}
