package k8smanage

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scaleClusterCmd = &cobra.Command{
	Use:   "scale cluster",
	Short: "Scale a Kubernetes cluster (provider-agnostic)",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		nodes, _ := cmd.Flags().GetInt("nodes")
		if provider == "" || name == "" || nodes < 1 {
			fmt.Println("Error: --provider, --name, and --nodes are required")
			cmd.Help()
			return
		}
		fmt.Printf("[stub] Would scale cluster '%s' with provider '%s' to %d nodes using shared abstraction.\n", name, provider, nodes)
		// TODO: Call shared/providers logic here
	},
}

func init() {
	scaleClusterCmd.Flags().String("provider", "", "Cloud provider (hetzner|aws|azure|gcp|...)")
	scaleClusterCmd.Flags().String("name", "", "Cluster name")
	scaleClusterCmd.Flags().Int("nodes", 1, "Number of nodes")
	RootCmd.AddCommand(scaleClusterCmd)
}
