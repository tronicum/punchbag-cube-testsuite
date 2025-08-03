package k8smanage

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createClusterCmd = &cobra.Command{
	Use:   "create cluster",
	Short: "Create a new Kubernetes cluster (provider-agnostic)",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		if provider == "" || name == "" {
			fmt.Println("Error: --provider and --name are required")
			cmd.Help()
			return
		}
		fmt.Printf("[stub] Would create cluster '%s' with provider '%s' using shared abstraction.\n", name, provider)
		// TODO: Call shared/providers logic here
	},
}

func init() {
	createClusterCmd.Flags().String("provider", "", "Cloud provider (hetzner|aws|azure|gcp|...)")
	createClusterCmd.Flags().String("name", "", "Cluster name")
	RootCmd.AddCommand(createClusterCmd)
}
