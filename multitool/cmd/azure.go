package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"punchbag-cube-testsuite/multitool/pkg/client"
	"punchbag-cube-testsuite/multitool/pkg/mock"
	"punchbag-cube-testsuite/multitool/pkg/models"
	"punchbag-cube-testsuite/multitool/pkg/output"

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
			output.Error(fmt.Sprintf("Failed to fetch Azure Monitor: %v", err))
			os.Exit(1)
		}
		defer resp.Body.Close()
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		f, err := os.Create(output)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to write file: %v", err))
			os.Exit(1)
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		enc.Encode(data)
		output.Success(fmt.Sprintf("Azure Monitor state saved to %s", output))
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
			output.Error(fmt.Sprintf("Failed to fetch Azure Log Analytics: %v", err))
			os.Exit(1)
		}
		defer resp.Body.Close()
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		f, err := os.Create(output)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to write file: %v", err))
			os.Exit(1)
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		enc.Encode(data)
		output.Success(fmt.Sprintf("Azure Log Analytics state saved to %s", output))
	},
}

// AKS management commands
var azureAksCmd = &cobra.Command{
	Use:   "aks",
	Short: "Manage Azure Kubernetes Service (AKS) clusters",
	Long: `Create, get, list, and delete AKS clusters in your Azure subscription.

Examples:
  multitool azure aks create --name my-aks --resource-group my-rg --location eastus
  multitool azure aks list
  multitool azure aks get --id <cluster-id>
  multitool azure aks delete --id <cluster-id>
`,
}

var azureAksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an AKS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		projectID, _ := cmd.Flags().GetString("project-id")
		name = promptIfEmpty(name, "AKS cluster name")
		resourceGroup = promptIfEmpty(resourceGroup, "Azure resource group name")
		location = promptIfEmpty(location, "Azure region/location")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		clusterClient := client.NewClusterClient(api)
		req := &models.ClusterCreateRequest{
			Name:          name,
			Provider:      models.Azure,
			ResourceGroup: resourceGroup,
			Location:      location,
			ProjectID:     projectID,
		}
		var cluster *models.Cluster
		var err error
		if simulateMode {
			cluster = mock.MockCreateAks(req)
		} else {
			cluster, err = clusterClient.CreateCluster(req)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create AKS cluster: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("AKS cluster created: %s (ID: %s)", cluster.Name, cluster.ID))
	},
}

var azureAksGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get AKS cluster details",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		clusterClient := client.NewClusterClient(api)
		var cluster *models.Cluster
		var err error
		if simulateMode {
			cluster = mock.MockGetAks(id)
		} else {
			cluster, err = clusterClient.GetCluster(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get AKS cluster: %v", err))
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
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		clusterClient := client.NewClusterClient(api)
		var clusters []*models.Cluster
		var err error
		if simulateMode {
			clusters = mock.MockListAks()
		} else {
			clusters, err = clusterClient.ListClustersByProvider(models.Azure)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list AKS clusters: %v", err))
			os.Exit(1)
		}
		if len(clusters) == 0 {
			output.Warn("No AKS clusters found.")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tResource Group\tLocation\tStatus")
		for _, c := range clusters {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", c.ID, c.Name, c.ResourceGroup, c.Location, c.Status)
		}
		w.Flush()
	},
}

var azureAksDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an AKS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		clusterClient := client.NewClusterClient(api)
		var err error
		if simulateMode {
			mock.MockDeleteAks(id)
		} else {
			err = clusterClient.DeleteCluster(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete AKS cluster: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("AKS cluster deleted: %s", id))
	},
}

// Log Analytics management commands
var azureLogCmd = &cobra.Command{
	Use:   "loganalytics",
	Short: "Manage Azure Log Analytics workspaces",
	Long: `Create, get, list, and delete Log Analytics workspaces in Azure.

Examples:
  multitool azure loganalytics create --name mylog --resource-group my-rg --location eastus --sku PerGB2018
  multitool azure loganalytics list
  multitool azure loganalytics get --id <workspace-id>
  multitool azure loganalytics delete --id <workspace-id>
`,
}

var azureLogCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Log Analytics workspace",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		sku, _ := cmd.Flags().GetString("sku")
		retention, _ := cmd.Flags().GetInt("retention-days")
		name = promptIfEmpty(name, "Log Analytics workspace name")
		resourceGroup = promptIfEmpty(resourceGroup, "Azure resource group name")
		location = promptIfEmpty(location, "Azure region/location")
		sku = promptIfEmpty(sku, "SKU (e.g. PerGB2018)")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		logClient := client.NewLogAnalyticsClient(api)
		workspace := &models.LogAnalyticsWorkspace{
			Name:          name,
			ResourceGroup: resourceGroup,
			Location:      location,
			Sku:           sku,
			RetentionDays: retention,
		}
		var result *models.LogAnalyticsWorkspace
		var err error
		if simulateMode {
			result = mock.MockCreateLogAnalytics(workspace)
		} else {
			result, err = logClient.Create(workspace)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create Log Analytics workspace: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Log Analytics workspace created: %s (ID: %s)", result.Name, result.ID))
	},
}

var azureLogGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Log Analytics workspace details",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		logClient := client.NewLogAnalyticsClient(api)
		var result *models.LogAnalyticsWorkspace
		var err error
		if simulateMode {
			result = mock.MockGetLogAnalytics(id)
		} else {
			result, err = logClient.Get(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get Log Analytics workspace: %v", err))
			os.Exit(1)
		}
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
	},
}

var azureLogListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Log Analytics workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		logClient := client.NewLogAnalyticsClient(api)
		var results []*models.LogAnalyticsWorkspace
		var err error
		if simulateMode {
			results = mock.MockListLogAnalytics()
		} else {
			results, err = logClient.List()
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list Log Analytics workspaces: %v", err))
			os.Exit(1)
		}
		if len(results) == 0 {
			output.Warn("No Log Analytics workspaces found.")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tResource Group\tLocation\tSKU\tRetention Days")
		for _, ws := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n", ws.ID, ws.Name, ws.ResourceGroup, ws.Location, ws.Sku, ws.RetentionDays)
		}
		w.Flush()
	},
}

var azureLogDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Log Analytics workspace",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		logClient := client.NewLogAnalyticsClient(api)
		var err error
		if simulateMode {
			mock.MockDeleteLogAnalytics(id)
		} else {
			err = logClient.Delete(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete Log Analytics workspace: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Log Analytics workspace deleted: %s", id))
	},
}

// Azure Budget management commands
var azureBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage Azure Budgets",
	Long: `Create, get, list, and delete Azure Budgets for cost management.

Examples:
  multitool azure budget create --name mybudget --resource-group my-rg --amount 100 --time-grain Monthly --start-date 2025-01-01 --end-date 2025-12-31
  multitool azure budget list
  multitool azure budget get --id <budget-id>
  multitool azure budget delete --id <budget-id>
`,
}

var azureBudgetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Azure Budget",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		amount, _ := cmd.Flags().GetFloat64("amount")
		timeGrain, _ := cmd.Flags().GetString("time-grain")
		startDate, _ := cmd.Flags().GetString("start-date")
		endDate, _ := cmd.Flags().GetString("end-date")
		name = promptIfEmpty(name, "Budget name")
		resourceGroup = promptIfEmpty(resourceGroup, "Azure resource group name")
		timeGrain = promptIfEmpty(timeGrain, "Time grain (e.g. Monthly)")
		startDate = promptIfEmpty(startDate, "Start date (YYYY-MM-DD)")
		endDate = promptIfEmpty(endDate, "End date (YYYY-MM-DD)")
		if amount == 0 {
			fmt.Print("Budget amount: ")
			fmt.Scanf("%f\n", &amount)
		}
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		budgetClient := client.NewAzureBudgetClient(api)
		budget := &models.AzureBudget{
			Name:          name,
			ResourceGroup: resourceGroup,
			Amount:        amount,
			TimeGrain:     timeGrain,
			StartDate:     startDate,
			EndDate:       endDate,
		}
		var result *models.AzureBudget
		var err error
		if simulateMode {
			result = mock.MockCreateBudget(budget)
		} else {
			result, err = budgetClient.Create(budget)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create Azure Budget: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Azure Budget created: %s (ID: %s)", result.Name, result.ID))
	},
}

var azureBudgetGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Azure Budget details",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		budgetClient := client.NewAzureBudgetClient(api)
		var result *models.AzureBudget
		var err error
		if simulateMode {
			result = mock.MockGetBudget(id)
		} else {
			result, err = budgetClient.Get(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get Azure Budget: %v", err))
			os.Exit(1)
		}
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
	},
}

var azureBudgetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Azure Budgets",
	Run: func(cmd *cobra.Command, args []string) {
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		budgetClient := client.NewAzureBudgetClient(api)
		var results []*models.AzureBudget
		var err error
		if simulateMode {
			results = mock.MockListBudgets()
		} else {
			results, err = budgetClient.List()
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list Azure Budgets: %v", err))
			os.Exit(1)
		}
		if len(results) == 0 {
			output.Warn("No Azure Budgets found.")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tResource Group\tAmount\tTime Grain\tStart Date\tEnd Date")
		for _, b := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%.2f\t%s\t%s\t%s\n", b.ID, b.Name, b.ResourceGroup, b.Amount, b.TimeGrain, b.StartDate, b.EndDate)
		}
		w.Flush()
	},
}

var azureBudgetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Azure Budget",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		budgetClient := client.NewAzureBudgetClient(api)
		var err error
		if simulateMode {
			mock.MockDeleteBudget(id)
		} else {
			err = budgetClient.Delete(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete Azure Budget: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Azure Budget deleted: %s", id))
	},
}

// Application Insights management commands
var azureAppInsightsCmd = &cobra.Command{
	Use:   "appinsights",
	Short: "Manage Azure Application Insights instances",
	Long: `Create, get, list, and delete Application Insights resources in Azure.

Examples:
  multitool azure appinsights create --name myapp --resource-group my-rg --location eastus --app-type web
  multitool azure appinsights list
  multitool azure appinsights get --id <appinsights-id>
  multitool azure appinsights delete --id <appinsights-id>
`,
}

var azureAppInsightsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Application Insights instance",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		appType, _ := cmd.Flags().GetString("app-type")
		retention, _ := cmd.Flags().GetInt("retention-days")
		name = promptIfEmpty(name, "App Insights name")
		resourceGroup = promptIfEmpty(resourceGroup, "Azure resource group name")
		location = promptIfEmpty(location, "Azure region/location")
		appType = promptIfEmpty(appType, "Application type (e.g. web)")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		appClient := client.NewAppInsightsClient(api)
		app := &models.AppInsightsResource{
			Name:          name,
			ResourceGroup: resourceGroup,
			Location:      location,
			AppType:       appType,
			RetentionDays: retention,
		}
		var result *models.AppInsightsResource
		var err error
		if simulateMode {
			result = mock.MockCreateAppInsights(app)
		} else {
			result, err = appClient.Create(app)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create Application Insights: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Application Insights created: %s (ID: %s)", result.Name, result.ID))
	},
}

var azureAppInsightsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Application Insights instance details",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		appClient := client.NewAppInsightsClient(api)
		var result *models.AppInsightsResource
		var err error
		if simulateMode {
			result = mock.MockGetAppInsights(id)
		} else {
			result, err = appClient.Get(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get Application Insights: %v", err))
			os.Exit(1)
		}
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
	},
}

var azureAppInsightsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Application Insights instances",
	Run: func(cmd *cobra.Command, args []string) {
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		appClient := client.NewAppInsightsClient(api)
		var results []*models.AppInsightsResource
		var err error
		if simulateMode {
			results = mock.MockListAppInsights()
		} else {
			results, err = appClient.List()
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list Application Insights: %v", err))
			os.Exit(1)
		}
		if len(results) == 0 {
			output.Warn("No Application Insights instances found.")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tResource Group\tLocation\tApp Type\tRetention Days")
		for _, a := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n", a.ID, a.Name, a.ResourceGroup, a.Location, a.AppType, a.RetentionDays)
		}
		w.Flush()
	},
}

var azureAppInsightsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Application Insights instance",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		api := client.NewAPIClient(proxyServer)
		if err := api.Authenticate(); err != nil {
			output.Error(fmt.Sprintf("Authentication failed (mocked): %v", err))
			os.Exit(1)
		}
		appClient := client.NewAppInsightsClient(api)
		var err error
		if simulateMode {
			mock.MockDeleteAppInsights(id)
		} else {
			err = appClient.Delete(id)
		}
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete Application Insights: %v", err))
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Application Insights deleted: %s", id))
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

	azureLogCreateCmd.Flags().String("name", "", "Log Analytics workspace Name")
	azureLogCreateCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureLogCreateCmd.Flags().String("location", "", "Azure region/location")
	azureLogCreateCmd.Flags().String("sku", "", "SKU (e.g. PerGB2018)")
	azureLogCreateCmd.Flags().Int("retention-days", 30, "Retention in days")
	azureLogCreateCmd.MarkFlagRequired("name")
	azureLogCreateCmd.MarkFlagRequired("resource-group")
	azureLogCreateCmd.MarkFlagRequired("location")
	azureLogCreateCmd.MarkFlagRequired("sku")

	azureLogGetCmd.Flags().String("id", "", "Log Analytics workspace ID")
	azureLogGetCmd.MarkFlagRequired("id")

	azureLogDeleteCmd.Flags().String("id", "", "Log Analytics workspace ID")
	azureLogDeleteCmd.MarkFlagRequired("id")

	azureBudgetCreateCmd.Flags().String("name", "", "Budget name")
	azureBudgetCreateCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureBudgetCreateCmd.Flags().Float64("amount", 0, "Budget amount")
	azureBudgetCreateCmd.Flags().String("time-grain", "", "Time grain (e.g. Monthly)")
	azureBudgetCreateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	azureBudgetCreateCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	azureBudgetCreateCmd.MarkFlagRequired("name")
	azureBudgetCreateCmd.MarkFlagRequired("resource-group")
	azureBudgetCreateCmd.MarkFlagRequired("amount")
	azureBudgetCreateCmd.MarkFlagRequired("time-grain")
	azureBudgetCreateCmd.MarkFlagRequired("start-date")
	azureBudgetCreateCmd.MarkFlagRequired("end-date")

	azureBudgetGetCmd.Flags().String("id", "", "Budget ID")
	azureBudgetGetCmd.MarkFlagRequired("id")

	azureBudgetDeleteCmd.Flags().String("id", "", "Budget ID")
	azureBudgetDeleteCmd.MarkFlagRequired("id")

	azureAppInsightsCreateCmd.Flags().String("name", "", "App Insights name")
	azureAppInsightsCreateCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureAppInsightsCreateCmd.Flags().String("location", "", "Azure region/location")
	azureAppInsightsCreateCmd.Flags().String("app-type", "", "Application type (e.g. web)")
	azureAppInsightsCreateCmd.Flags().Int("retention-days", 90, "Retention in days")
	azureAppInsightsCreateCmd.MarkFlagRequired("name")
	azureAppInsightsCreateCmd.MarkFlagRequired("resource-group")
	azureAppInsightsCreateCmd.MarkFlagRequired("location")
	azureAppInsightsCreateCmd.MarkFlagRequired("app-type")

	azureAppInsightsGetCmd.Flags().String("id", "", "App Insights ID")
	azureAppInsightsGetCmd.MarkFlagRequired("id")

	azureAppInsightsDeleteCmd.Flags().String("id", "", "App Insights ID")
	azureAppInsightsDeleteCmd.MarkFlagRequired("id")
}

func promptIfEmpty(val, prompt string) string {
	if val != "" {
		return val
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
