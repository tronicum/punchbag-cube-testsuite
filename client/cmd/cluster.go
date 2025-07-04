package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage clusters across multiple cloud providers",
	Long:  `Commands for managing clusters in the test suite across Azure, StackIT, AWS, and GCP.`,
}

var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))
		
		provider, _ := cmd.Flags().GetString("provider")
		var clusters []*api.Cluster
		var err error
		
		if provider != "" {
			clusters, err = client.ListClustersByProvider(provider)
		} else {
			clusters, err = client.ListClusters()
		}
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing clusters: %v\n", err)
			os.Exit(1)
		}

		output.PrintClusters(clusters, viper.GetString("format"))
	},
}

var clusterGetCmd = &cobra.Command{
	Use:   "get [cluster-id]",
	Short: "Get cluster details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))
		cluster, err := client.GetCluster(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting cluster: %v\n", err)
			os.Exit(1)
		}

		output.PrintCluster(cluster, viper.GetString("format"))
	},
}

var clusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cluster",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		provider, _ := cmd.Flags().GetString("provider")
		
		if name == "" || provider == "" {
			fmt.Fprintf(os.Stderr, "Error: name and provider are required\n")
			os.Exit(1)
		}

		client := api.NewClient(viper.GetString("server"))
		var created *api.Cluster
		var err error
		
		switch provider {
		case "azure":
			resourceGroup, _ := cmd.Flags().GetString("resource-group")
			location, _ := cmd.Flags().GetString("location")
			
			if resourceGroup == "" || location == "" {
				fmt.Fprintf(os.Stderr, "Error: resource-group and location are required for Azure clusters\n")
				os.Exit(1)
			}
			
			config := make(map[string]interface{})
			if kubernetesVersion, _ := cmd.Flags().GetString("kubernetes-version"); kubernetesVersion != "" {
				config["kubernetes_version"] = kubernetesVersion
			}
			if nodeCount, _ := cmd.Flags().GetInt("node-count"); nodeCount > 0 {
				config["node_count"] = nodeCount
			}
			
			created, err = client.CreateAzureCluster(name, resourceGroup, location, config)
			
		case "schwarz-stackit":
			projectID, _ := cmd.Flags().GetString("project-id")
			region, _ := cmd.Flags().GetString("region")
			
			if projectID == "" || region == "" {
				fmt.Fprintf(os.Stderr, "Error: project-id and region are required for StackIT clusters\n")
				os.Exit(1)
			}
			
			config := make(map[string]interface{})
			if kubernetesVersion, _ := cmd.Flags().GetString("kubernetes-version"); kubernetesVersion != "" {
				config["kubernetes_version"] = kubernetesVersion
			}
			if nodeCount, _ := cmd.Flags().GetInt("node-count"); nodeCount > 0 {
				config["node_count"] = nodeCount
			}
			
			created, err = client.CreateStackITCluster(name, projectID, region, config)
			
		default:
			fmt.Fprintf(os.Stderr, "Error: unsupported provider '%s'. Supported providers: azure, schwarz-stackit\n", provider)
			os.Exit(1)
		}
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating cluster: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cluster created successfully with ID: %s\n", created.ID)
		output.PrintCluster(created, viper.GetString("format"))
	},
}

var clusterDeleteCmd = &cobra.Command{
	Use:   "delete [cluster-id]",
	Short: "Delete a cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		confirm, _ := cmd.Flags().GetBool("confirm")
		if !confirm {
			fmt.Print("Are you sure you want to delete this cluster? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Operation cancelled")
				return
			}
		}

		client := api.NewClient(viper.GetString("server"))
		err := client.DeleteCluster(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting cluster: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cluster %s deleted successfully\n", args[0])
	},
}

var clusterTestCmd = &cobra.Command{
	Use:   "test [cluster-id]",
	Short: "Run a test on a cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		testType, _ := cmd.Flags().GetString("type")
		configFile, _ := cmd.Flags().GetString("config")

		if testType == "" {
			fmt.Fprintf(os.Stderr, "Error: test type is required\n")
			os.Exit(1)
		}

		config := make(map[string]interface{})
		if configFile != "" {
			data, err := os.ReadFile(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
				os.Exit(1)
			}
			if err := json.Unmarshal(data, &config); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
				os.Exit(1)
			}
		}

		testReq := &api.TestRequest{
			ClusterID: args[0],
			TestType:  testType,
			Config:    config,
		}

		client := api.NewClient(viper.GetString("server"))
		result, err := client.RunTest(args[0], testReq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running test: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Test started with ID: %s\n", result.ID)
		output.PrintTestResult(result, viper.GetString("format"))
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterListCmd)
	clusterCmd.AddCommand(clusterGetCmd)
	clusterCmd.AddCommand(clusterCreateCmd)
	clusterCmd.AddCommand(clusterDeleteCmd)
	clusterCmd.AddCommand(clusterTestCmd)

	// List cluster flags
	clusterListCmd.Flags().String("provider", "", "Filter by cloud provider (azure, schwarz-stackit, aws, gcp)")

	// Create cluster flags
	clusterCreateCmd.Flags().String("name", "", "Cluster name (required)")
	clusterCreateCmd.Flags().String("provider", "", "Cloud provider (required: azure, schwarz-stackit)")
	
	// Azure-specific flags
	clusterCreateCmd.Flags().String("resource-group", "", "Azure resource group (required for Azure)")
	clusterCreateCmd.Flags().String("location", "", "Azure location (required for Azure)")
	
	// StackIT-specific flags
	clusterCreateCmd.Flags().String("project-id", "", "StackIT project ID (required for StackIT)")
	clusterCreateCmd.Flags().String("region", "", "StackIT region (required for StackIT)")
	
	// Common flags
	clusterCreateCmd.Flags().String("kubernetes-version", "", "Kubernetes version")
	clusterCreateCmd.Flags().Int("node-count", 0, "Number of nodes")

	// Delete cluster flags
	clusterDeleteCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")

	// Test cluster flags
	clusterTestCmd.Flags().String("type", "", "Test type (load_test, performance_test, stress_test)")
	clusterTestCmd.Flags().String("config", "", "Test configuration file (JSON)")
}
