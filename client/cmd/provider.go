package cmd

import (
	"fmt"

	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"

	"github.com/spf13/cobra"
)

// providerCmd represents the provider command
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage cloud providers and their operations",
	Long: `Manage cloud providers and execute provider-specific operations.
Supported providers: azure, hetzner-hcloud, united-ionos, schwarz-stackit, aws, gcp`,
}

// providerInfoCmd gets information about a provider
var providerInfoCmd = &cobra.Command{
	Use:   "info [provider]",
	Short: "Get information about a cloud provider",
	Long: `Get detailed information about a cloud provider including features, 
documentation links, and pricing model.`,
	Args: cobra.ExactArgs(1),
	RunE: runProviderInfo,
}

// providerListCmd lists clusters for a specific provider
var providerListCmd = &cobra.Command{
	Use:   "list [provider]",
	Short: "List clusters for a specific cloud provider",
	Long:  `List all clusters managed by a specific cloud provider.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProviderList,
}

// providerOperationCmd executes a provider-specific operation
var providerOperationCmd = &cobra.Command{
	Use:   "operation [provider] [operation]",
	Short: "Execute a provider-specific operation",
	Long: `Execute a provider-specific operation such as create, delete, scale, etc.
The operation parameters can be provided as JSON via the --params flag.`,
	Args: cobra.ExactArgs(2),
	RunE: runProviderOperation,
}

var operationParams string

func init() {
	rootCmd.AddCommand(providerCmd)
	providerCmd.AddCommand(providerInfoCmd)
	providerCmd.AddCommand(providerListCmd)
	providerCmd.AddCommand(providerOperationCmd)
	
	providerOperationCmd.Flags().StringVar(&operationParams, "params", "{}", "Operation parameters as JSON")
}

func runProviderInfo(cmd *cobra.Command, args []string) error {
	provider := args[0]
	
	client := api.NewClient(apiBaseURL)
	
	info, err := client.GetProviderInfo(provider)
	if err != nil {
		return fmt.Errorf("failed to get provider info for %s: %w", provider, err)
	}
	
	formatter := output.NewFormatter(outputFormat)
	return formatter.FormatProviderInfo(info)
}

func runProviderList(cmd *cobra.Command, args []string) error {
	provider := args[0]
	
	client := api.NewClient(apiBaseURL)
	
	clusters, err := client.ListProviderClusters(provider)
	if err != nil {
		return fmt.Errorf("failed to list clusters for provider %s: %w", provider, err)
	}
	
	formatter := output.NewFormatter(outputFormat)
	return formatter.FormatProviderClusters(clusters)
}

func runProviderOperation(cmd *cobra.Command, args []string) error {
	provider := args[0]
	operation := args[1]
	
	client := api.NewClient(apiBaseURL)
	
	result, err := client.ExecuteProviderOperation(provider, operation, operationParams)
	if err != nil {
		return fmt.Errorf("failed to execute %s operation for provider %s: %w", operation, provider, err)
	}
	
	formatter := output.NewFormatter(outputFormat)
	return formatter.FormatProviderOperation(result)
}
