package cmd

import (
	"github.com/spf13/cobra"
)

// k8sManageCmd is the top-level cluster lifecycle management command
var k8sManageCmd = &cobra.Command{
	Use:   "k8s-manage",
	Short: "Cluster lifecycle management (create, delete, scale, upgrade, etc.)",
}

func init() {
	rootCmd.AddCommand(k8sManageCmd)
}
