package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/client"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/output"
	importpkg "github.com/tronicum/punchbag-cube-testsuite/shared/import"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

var (
	serverURL     string
	outputFormat  string
	resourceGroup string
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
		clusterName := args[0]
		providerStr := args[1]

		provider := sharedmodels.CloudProvider(providerStr)
		if !isValidProvider(provider) {
			output.FormatError(fmt.Errorf("invalid provider: %s. Supported providers: azure, aws, gcp, hetzner, ionos, stackit", providerStr))
			os.Exit(1)
		}

		apiClient := client.NewAPIClient(serverURL)
		clusterClient := client.NewClusterClient(apiClient)

		// Build provider config
		providerConfig := make(map[string]interface{})
		config := make(map[string]interface{})

		switch provider {
		case sharedmodels.Azure:
			if resourceGroup == "" {
				output.FormatError(fmt.Errorf("resource-group is required for Azure"))
				os.Exit(1)
			}
			if location == "" {
				output.FormatError(fmt.Errorf("location is required for Azure"))
				os.Exit(1)
			}
			providerConfig["resource_group"] = resourceGroup
			providerConfig["location"] = location
		case sharedmodels.AWS:
			if region == "" {
				output.FormatError(fmt.Errorf("region is required for AWS"))
				os.Exit(1)
			}
			providerConfig["region"] = region
		case sharedmodels.GCP:
			if projectID == "" {
				output.FormatError(fmt.Errorf("project-id is required for GCP"))
				os.Exit(1)
			}
			if region == "" {
				output.FormatError(fmt.Errorf("region is required for GCP"))
				os.Exit(1)
			}
			providerConfig["project_id"] = projectID
			providerConfig["region"] = region
		}

		// Load additional config from file if provided
		if configFile != "" {
			f, err := os.Open(configFile)
			if err != nil {
				output.FormatError(fmt.Errorf("failed to open config file: %w", err))
				os.Exit(1)
			}
			defer f.Close()
			_, err = importpkg.LoadConfigJSON(f)
			if err != nil {
				output.FormatError(fmt.Errorf("failed to load config file: %w", err))
				os.Exit(1)
			}
			// TODO: Merge loaded config into config map (extend as needed)
		}

		req := &sharedmodels.ClusterCreateRequest{
			Name:           clusterName,
			Provider:       provider,
			Config:         config,
			ProviderConfig: providerConfig,
			ProjectID:      projectID,
			ResourceGroup:  resourceGroup,
			Location:       location,
			Region:         region,
		}

		output.FormatInfo(fmt.Sprintf("Creating cluster '%s' on %s...", clusterName, provider))

		cluster, err := clusterClient.CreateCluster(req)
		if err != nil {
			output.FormatError(fmt.Errorf("failed to create cluster: %w", err))
			os.Exit(1)
		}

		output.FormatSuccess(fmt.Sprintf("Cluster '%s' created successfully with ID: %s", cluster.Name, cluster.ID))

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(cluster); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
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
		apiClient := client.NewAPIClient(serverURL)
		clusterClient := client.NewClusterClient(apiClient)

		var clusters []*sharedmodels.Cluster
		var err error

		if len(args) > 0 {
			provider := sharedmodels.CloudProvider(args[0])
			if !isValidProvider(provider) {
				output.FormatError(fmt.Errorf("invalid provider: %s", args[0]))
				os.Exit(1)
			}
			clusters, err = clusterClient.ListClustersByProvider(provider)
		} else {
			clusters, err = clusterClient.ListClusters()
		}

		if err != nil {
			output.FormatError(fmt.Errorf("failed to list clusters: %w", err))
			os.Exit(1)
		}

		if len(clusters) == 0 {
			output.FormatInfo("No clusters found")
			return
		}

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(clusters); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
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
		clusterID := args[0]

		apiClient := client.NewAPIClient(serverURL)
		clusterClient := client.NewClusterClient(apiClient)

		cluster, err := clusterClient.GetCluster(clusterID)
		if err != nil {
			output.FormatError(fmt.Errorf("failed to get cluster: %w", err))
			os.Exit(1)
		}

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(cluster); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
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
		clusterID := args[0]
		confirm, _ := cmd.Flags().GetBool("confirm")

		if !confirm {
			fmt.Printf("Are you sure you want to delete cluster '%s'? This action cannot be undone. (y/N): ", clusterID)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				output.FormatInfo("Deletion cancelled")
				return
			}
		}

		apiClient := client.NewAPIClient(serverURL)
		clusterClient := client.NewClusterClient(apiClient)

		output.FormatInfo(fmt.Sprintf("Deleting cluster '%s'...", clusterID))

		err := clusterClient.DeleteCluster(clusterID)
		if err != nil {
			output.FormatError(fmt.Errorf("failed to delete cluster: %w", err))
			os.Exit(1)
		}

		output.FormatSuccess(fmt.Sprintf("Cluster '%s' deleted successfully", clusterID))
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
		clusterID := args[0]
		testType := args[1]

		apiClient := client.NewAPIClient(serverURL)
		testClient := client.NewTestClient(apiClient)

		// Load test config from file if provided
		testConfig := make(map[string]interface{})
		if configFile != "" {
			f, err := os.Open(configFile)
			if err != nil {
				output.FormatError(fmt.Errorf("failed to open config file: %w", err))
				os.Exit(1)
			}
			defer f.Close()
			_, err = importpkg.LoadConfigJSON(f)
			if err != nil {
				output.FormatError(fmt.Errorf("failed to load config file: %w", err))
				os.Exit(1)
			}
			// TODO: Map config fields as needed
		}

		req := &sharedmodels.TestRequest{
			ClusterID: clusterID,
			TestType:  testType,
			Config:    testConfig,
		}

		output.FormatInfo(fmt.Sprintf("Starting %s test on cluster '%s'...", testType, clusterID))

		result, err := testClient.RunTest(req)
		if err != nil {
			output.FormatError(fmt.Errorf("failed to run test: %w", err))
			os.Exit(1)
		}

		output.FormatSuccess(fmt.Sprintf("Test started with ID: %s", result.ID))

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(result); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
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
		clusterID := args[0]

		apiClient := client.NewAPIClient(serverURL)
		testClient := client.NewTestClient(apiClient)

		results, err := testClient.ListTestResults(clusterID)
		if err != nil {
			output.FormatError(fmt.Errorf("failed to list test results: %w", err))
			os.Exit(1)
		}

		if len(results) == 0 {
			output.FormatInfo("No test results found")
			return
		}

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(results); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
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
		testID := args[0]

		apiClient := client.NewAPIClient(serverURL)
		testClient := client.NewTestClient(apiClient)

		result, err := testClient.GetTestResult(testID)
		if err != nil {
			output.FormatError(fmt.Errorf("failed to get test result: %w", err))
			os.Exit(1)
		}

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(result); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
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
func httpPost(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

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
	clusterCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "Server URL")
	clusterCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	testCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "Server URL")
	testCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")

	// Cluster create flags
	clusterCreateCmd.Flags().StringVar(&resourceGroup, "resource-group", "", "Azure resource group")
	clusterCreateCmd.Flags().StringVar(&location, "location", "", "Azure location")
	clusterCreateCmd.Flags().StringVar(&region, "region", "", "AWS/GCP region")
	clusterCreateCmd.Flags().StringVar(&projectID, "project-id", "", "GCP project ID")
	clusterCreateCmd.Flags().StringVar(&configFile, "config", "", "Configuration file (JSON)")

	// Cluster delete flags
	clusterDeleteCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")

	// Test run flags
	testRunCmd.Flags().StringVar(&configFile, "config", "", "Test configuration file (JSON)")
}
