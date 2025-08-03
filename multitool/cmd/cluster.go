package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/output"
	"github.com/tronicum/punchbag-cube-testsuite/shared/log"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

var (
outputFormat  string
	location      string
	region        string
	projectID     string
	configFile    string
)

// clusterCmd represents the cluster command group
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage Kubernetes clusters across cloud providers",
	Long:  "Create, list, get, and delete Kubernetes clusters on Azure, AWS, GCP, and other cloud providers.",
}

// clusterCreateCmd creates a new cluster
var clusterCreateCmd = &cobra.Command{
	Use:   "create [name] [provider]",
	Short: "Create a new Kubernetes cluster",
	Long: `Create a new Kubernetes cluster on the specified cloud provider.
	
Supported providers: azure, aws, gcp, hetzner, ionos, stackit

Examples:
  multitool cluster create my-cluster azure --resource-group my-rg --location eastus
  multitool cluster create my-cluster aws --region us-west-2
  multitool cluster create my-cluster gcp --project-id my-project --region us-central1`,
	Args: cobra.ExactArgs(2),
Run: func(cmd *cobra.Command, args []string) {
	// clusterName := args[0]
	providerStr := args[1]

	// Read output format from flag, default to 'table' if not set
	outputFormatFlag, err := cmd.Flags().GetString("output")
	if err != nil || outputFormatFlag == "" {
		outputFormatFlag = "table"
	}

	provider := sharedmodels.CloudProvider(providerStr)
	if !isValidProvider(provider) {
		err := errors.New("invalid provider: " + providerStr + ". Supported providers: azure, aws, gcp, hetzner, ionos, stackit")
		output.FormatError(err)
		log.Error(err.Error())
		return
	}

	// TODO: Use shared library for cluster operations (create)
	// Use outputFormatFlag for output formatting
	// ...existing code for providerConfig and config...
	// ...existing code for loading configFile...
	// TODO: Call shared library to create cluster and print result
},
}

// clusterListCmd lists clusters
var clusterListCmd = &cobra.Command{
	Use:   "list [provider]",
	Short: "List Kubernetes clusters",
	Long: `List all Kubernetes clusters or filter by cloud provider.
	
Examples:
  multitool cluster list
  multitool cluster list azure
  multitool cluster list aws`,
	Args: cobra.MaximumNArgs(1),
Run: func(cmd *cobra.Command, args []string) {
	outputFormatFlag, err := cmd.Flags().GetString("output")
	if err != nil || outputFormatFlag == "" {
		outputFormatFlag = "table"
	}
	// TODO: Call shared library to list clusters and print result using outputFormatFlag
},
}

// clusterGetCmd gets a specific cluster
var clusterGetCmd = &cobra.Command{
	Use:   "get [cluster-id]",
	Short: "Get details of a specific cluster",
	Long: `Get detailed information about a specific Kubernetes cluster.
	
Examples:
  multitool cluster get cluster-123
  multitool cluster get cluster-123 --output json`,
	Args: cobra.ExactArgs(1),
Run: func(cmd *cobra.Command, args []string) {
	outputFormatFlag, err := cmd.Flags().GetString("output")
	if err != nil || outputFormatFlag == "" {
		outputFormatFlag = "table"
	}
	// clusterID := args[0]
	// TODO: Call shared library to get cluster and print result using outputFormatFlag
},
}

// clusterDeleteCmd deletes a cluster
var clusterDeleteCmd = &cobra.Command{
	Use:   "delete [cluster-id]",
	Short: "Delete a Kubernetes cluster",
	Long: `Delete a Kubernetes cluster. This operation is irreversible.
	
Examples:
  multitool cluster delete cluster-123
  multitool cluster delete cluster-123 --confirm`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// clusterID := args[0]
		confirm, _ := cmd.Flags().GetBool("confirm")

		if !confirm {
			log.Info("Are you sure you want to delete this cluster? This action cannot be undone. (y/N): ")
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				output.FormatError(err)
				log.Error("failed to read confirmation input: %v", err)
				return
			}
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				output.FormatInfo("Deletion cancelled")
				log.Info("Deletion cancelled")
				return
			}
		}

		// TODO: Call shared library to delete cluster and print result
	},
}

// testCmd represents the test command group
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Manage cluster tests",
	Long:  "Run and manage tests on Kubernetes clusters.",
}

// testRunCmd runs a test on a cluster
var testRunCmd = &cobra.Command{
	Use:   "run [cluster-id] [test-type]",
	Short: "Run a test on a cluster",
	Long: `Run a specific test on a Kubernetes cluster.
	
Supported test types: connectivity, performance, security, compliance
	
Examples:
  multitool test run cluster-123 connectivity
  multitool test run cluster-123 performance --config perf-config.json`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// clusterID := args[0]
		// testType := args[1]

		// TODO: Call shared library to run test and print result
	},
}

// testListCmd lists test results
var testListCmd = &cobra.Command{
	Use:   "list [cluster-id]",
	Short: "List test results for a cluster",
	Long: `List all test results for a specific cluster.
	
Examples:
  multitool test list cluster-123
  multitool test list cluster-123 --output json`,
	Args: cobra.ExactArgs(1),
Run: func(cmd *cobra.Command, args []string) {
	outputFormatFlag, err := cmd.Flags().GetString("output")
	if err != nil || outputFormatFlag == "" {
		outputFormatFlag = "table"
	}
	// clusterID := args[0]
	// TODO: Call shared library to list test results and print result using outputFormatFlag
},
}

// testGetCmd gets a specific test result
var testGetCmd = &cobra.Command{
	Use:   "get [test-id]",
	Short: "Get details of a specific test result",
	Long: `Get detailed information about a specific test result.
	
Examples:
  multitool test get test-456
  multitool test get test-456 --output yaml`,
	Args: cobra.ExactArgs(1),
Run: func(cmd *cobra.Command, args []string) {
	outputFormatFlag, err := cmd.Flags().GetString("output")
	if err != nil || outputFormatFlag == "" {
		outputFormatFlag = "table"
	}
	// testID := args[0]
	// TODO: Call shared library to get test result and print result using outputFormatFlag
},
}

// Helper functions

func isValidProvider(provider sharedmodels.CloudProvider) bool {
	validProviders := []sharedmodels.CloudProvider{
		sharedmodels.CloudProviderAzure,
		sharedmodels.CloudProviderAWS,
		sharedmodels.CloudProviderGCP,
		sharedmodels.CloudProviderHetzner,
		sharedmodels.CloudProviderIONOS,
		sharedmodels.CloudProviderStackIT,
	}
	for _, p := range validProviders {
		if provider == p {
			return true
		}
	}
	return false
}

// Config loading is now handled by shared/import. See above for usage.

// Helper for HTTP POST in proxy mode

func init() {
	// Add cluster subcommands
	clusterCmd.AddCommand(clusterCreateCmd)
	clusterCmd.AddCommand(clusterListCmd)
	clusterCmd.AddCommand(clusterGetCmd)
	clusterCmd.AddCommand(clusterDeleteCmd)

	// Add test subcommands
	testCmd.AddCommand(testRunCmd)
	testCmd.AddCommand(testListCmd)
	testCmd.AddCommand(testGetCmd)

	// Global flags
	clusterCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	testCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	clusterCreateCmd.Flags().StringVar(&location, "location", "", "Azure location")
	clusterCreateCmd.Flags().StringVar(&region, "region", "", "AWS/GCP region")
	clusterCreateCmd.Flags().StringVar(&projectID, "project-id", "", "GCP project ID")
	clusterCreateCmd.Flags().StringVar(&configFile, "config", "", "Configuration file (JSON)")

	// Cluster delete flags
	clusterDeleteCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")

	// Test run flags
	testRunCmd.Flags().StringVar(&configFile, "config", "", "Test configuration file (JSON)")
}
