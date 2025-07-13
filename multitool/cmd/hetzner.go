package cmd

import (
	"github.com/spf13/cobra"
)

var hetznerCmd = &cobra.Command{
	Use:   "hetzner",
	Short: "Hetzner cloud provider operations",
	Long:  `Manage Hetzner resources (placeholder, extend with subcommands).`,
}

func init() {
	// Add Hetzner subcommands here
}
