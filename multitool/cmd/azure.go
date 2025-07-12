package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	azureprovider "github.com/tronicum/punchbag-cube-testsuite/shared/providers/azure"
)

// Root Azure Command
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Azure cloud provider operations",
	Long:  `Manage Azure resources including AKS clusters, monitoring, budgets, and storage.`,
}

// ==== MONITORING ====

var azureMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Manage Azure Monitor resources",
	Long:  `Create, update, and manage Azure Monitor services.`,
}

var azureCreateMonitorCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Azure Monitor resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		workspaceName, _ := cmd.Flags().GetString("workspace-name")
		simulationMode, _ := cmd.Flags().GetBool("simulation")

		provider := azureprovider.NewAzureProvider()
		provider.SetSimulationMode(simulationMode)

		ctx := context.Background()
		result, err := provider.CreateMonitor(ctx, resourceGroup, location, workspaceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Azure Monitor: %v\n", err)
			return err
		}
		fmt.Printf("Azure Monitor created: ID=%s, Status=%s, Resources=%v\n", result.ID, result.Status, result.Resources)
		return nil
	},
}

// ==== BUDGET ====

var azureBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage Azure Budget resources",
	Long:  `Create, update, and manage Azure Budget resources.`,
}

var azureCreateBudgetCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Azure Budget for cost management",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		amount, _ := cmd.Flags().GetFloat64("amount")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		timeGrain, _ := cmd.Flags().GetString("time-grain")
		simulationMode, _ := cmd.Flags().GetBool("simulation")

		provider := azureprovider.NewAzureProvider()
		provider.SetSimulationMode(simulationMode)

		ctx := context.Background()
		result, err := provider.CreateBudget(ctx, name, amount, resourceGroup, timeGrain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Azure Budget: %v\n", err)
			return err
		}
		fmt.Printf("Azure Budget created: ID=%s, Status=%s\n", result.ID, result.Status)
		return nil
	},
}

// ==== AKS ====

var azureAksCmd = &cobra.Command{
	Use:   "aks",
	Short: "Manage Azure Kubernetes Service clusters",
	Long:  `Create, update, scale, and manage Azure AKS clusters.`,
}

var azureCreateAksCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Azure Kubernetes Service cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		nodeCount, _ := cmd.Flags().GetInt("node-count")
		simulationMode, _ := cmd.Flags().GetBool("simulation")

		provider := azureprovider.NewAzureProvider()
		provider.SetSimulationMode(simulationMode)

		ctx := context.Background()
		result, err := provider.CreateAKSCluster(ctx, name, resourceGroup, location, nodeCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating AKS cluster: %v\n", err)
			return err
		}
		fmt.Printf("AKS Cluster created: ID=%s, Status=%s, URL=%s\n", result.ID, result.Status, result.URL)
		return nil
	},
}

// ==== MONITORING STACK ====
// (Optional: Implement stack commands using provider if needed)

// ==== COMMAND TREE & FLAGS ====

func init() {
	rootCmd.AddCommand(azureCmd)

	// Azure subcommands
	azureCmd.AddCommand(azureMonitorCmd)
	azureCmd.AddCommand(azureBudgetCmd)
	azureCmd.AddCommand(azureAksCmd)

	// Azure Monitor
	azureMonitorCmd.AddCommand(azureCreateMonitorCmd)
	azureCreateMonitorCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateMonitorCmd.Flags().String("location", "eastus", "Azure region")
	azureCreateMonitorCmd.Flags().String("workspace-name", "", "Log Analytics workspace name")
	azureCreateMonitorCmd.Flags().Bool("simulation", false, "Use simulation mode")
	azureCreateMonitorCmd.MarkFlagRequired("resource-group")
	azureCreateMonitorCmd.MarkFlagRequired("workspace-name")

	// Azure Budget
	azureBudgetCmd.AddCommand(azureCreateBudgetCmd)
	azureCreateBudgetCmd.Flags().String("name", "", "Budget name")
	azureCreateBudgetCmd.Flags().Float64("amount", 0, "Budget amount in USD")
	azureCreateBudgetCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateBudgetCmd.Flags().String("time-grain", "Monthly", "Budget time grain")
	azureCreateBudgetCmd.Flags().Bool("simulation", false, "Use simulation mode")
	azureCreateBudgetCmd.MarkFlagRequired("name")
	azureCreateBudgetCmd.MarkFlagRequired("amount")
	azureCreateBudgetCmd.MarkFlagRequired("resource-group")

	// Azure AKS
	azureAksCmd.AddCommand(azureCreateAksCmd)
	azureCreateAksCmd.Flags().String("name", "", "AKS cluster name")
	azureCreateAksCmd.Flags().String("resource-group", "", "Azure resource group name")
	azureCreateAksCmd.Flags().String("location", "eastus", "Azure region")
	azureCreateAksCmd.Flags().Int("node-count", 3, "Number of nodes in default pool")
	azureCreateAksCmd.Flags().Bool("simulation", false, "Use simulation mode")
	azureCreateAksCmd.MarkFlagRequired("name")
	azureCreateAksCmd.MarkFlagRequired("resource-group")
}
