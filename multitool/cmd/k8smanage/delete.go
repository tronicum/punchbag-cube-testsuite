package k8smanage

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteClusterCmd = &cobra.Command{
	Use:   "delete cluster",
	Short: "Delete a Kubernetes cluster (provider-agnostic)",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		if provider == "" || name == "" {
			fmt.Println("Error: --provider and --name are required")
			cmd.Help()
			return
		}
		fmt.Printf("[stub] Would delete cluster '%s' with provider '%s' using shared abstraction.\n", name, provider)
		// TODO: Call shared/providers logic here
	},
}

func init() {
	deleteClusterCmd.Flags().String("provider", "", "Cloud provider (hetzner|aws|azure|gcp|...)")
	deleteClusterCmd.Flags().String("name", "", "Cluster name")
	RootCmd.AddCommand(deleteClusterCmd)
}
