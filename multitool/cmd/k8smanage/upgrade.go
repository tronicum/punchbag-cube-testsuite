package k8smanage

import (
	"fmt"
	"github.com/spf13/cobra"
)

var upgradeClusterCmd = &cobra.Command{
	Use:   "upgrade cluster",
	Short: "Upgrade a Kubernetes cluster (provider-agnostic)",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		if provider == "" || name == "" || version == "" {
			fmt.Println("Error: --provider, --name, and --version are required")
			cmd.Help()
			return
		}
		fmt.Printf("[stub] Would upgrade cluster '%s' with provider '%s' to version '%s' using shared abstraction.\n", name, provider, version)
		// TODO: Call shared/providers logic here
	},
}

func init() {
	upgradeClusterCmd.Flags().String("provider", "", "Cloud provider (hetzner|aws|azure|gcp|...)")
	upgradeClusterCmd.Flags().String("name", "", "Cluster name")
	upgradeClusterCmd.Flags().String("version", "", "Kubernetes version")
	RootCmd.AddCommand(upgradeClusterCmd)
}
