package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"punchbag-cube-testsuite/client/pkg/api"
	"punchbag-cube-testsuite/client/pkg/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage AKS clusters",
	Long:  `Commands for managing AKS clusters in the test suite.`,
}

var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(viper.GetString("server"))
		clusters, err := client.ListClusters()
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
		resourceGroup, _ := cmd.Flags().GetString("resource-group")
		location, _ := cmd.Flags().GetString("location")
		kubernetesVersion, _ := cmd.Flags().GetString("kubernetes-version")
		nodeCount, _ := cmd.Flags().GetInt("node-count")

		if name == "" || resourceGroup == "" || location == "" {
			fmt.Fprintf(os.Stderr, "Error: name, resource-group, and location are required\n")
			os.Exit(1)
		}

		cluster := &api.AKSCluster{
			Name:              name,
			ResourceGroup:     resourceGroup,
			Location:          location,
			KubernetesVersion: kubernetesVersion,
			NodeCount:         nodeCount,
			Status:           "creating",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		client := api.NewClient(viper.GetString("server"))
		created, err := client.CreateCluster(cluster)
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

		testReq := &api.AKSTestRequest{
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

	// Create cluster flags
	clusterCreateCmd.Flags().String("name", "", "Cluster name (required)")
	clusterCreateCmd.Flags().String("resource-group", "", "Azure resource group (required)")
	clusterCreateCmd.Flags().String("location", "", "Azure location (required)")
	clusterCreateCmd.Flags().String("kubernetes-version", "1.28.0", "Kubernetes version")
	clusterCreateCmd.Flags().Int("node-count", 3, "Number of nodes")

	// Delete cluster flags
	clusterDeleteCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")

	// Test cluster flags
	clusterTestCmd.Flags().String("type", "", "Test type (load_test, performance_test, stress_test)")
	clusterTestCmd.Flags().String("config", "", "Test configuration file (JSON)")
}
