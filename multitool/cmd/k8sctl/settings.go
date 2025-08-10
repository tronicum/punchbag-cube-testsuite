// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl settings: kubectl-like settings command for multitool (future extension)
package k8sctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var settingsCmd = &cobra.Command{
	Use:   "settings [flags]",
	Short: "Manage k8sctl settings (future extension)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("k8sctl settings: not yet implemented. This will manage CLI/user settings in the future.")
	},
}

func init() {
	RootCmd.AddCommand(settingsCmd)
}
