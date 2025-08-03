// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl plugin: manage kubectl plugins via krew
package k8sctl

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage kubectl plugins via krew (install, list, upgrade, etc.)",
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [plugin]",
	Short: "Install a kubectl plugin via krew",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plugin := args[0]
		fmt.Printf("Installing kubectl plugin '%s' via krew...\n", plugin)
		out, err := exec.Command("kubectl", "krew", "install", plugin).CombinedOutput()
		if err != nil {
			fmt.Printf("Error installing plugin: %v\n%s\n", err, string(out))
			return
		}
		fmt.Println("Plugin installed successfully.")
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed kubectl plugins via krew",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("kubectl", "krew", "list").CombinedOutput()
		if err != nil {
			fmt.Printf("Error listing plugins: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

var pluginUpgradeCmd = &cobra.Command{
	Use:   "upgrade [plugin]",
	Short: "Upgrade a kubectl plugin via krew (or all if omitted)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var out []byte
		var err error
		if len(args) == 0 {
			out, err = exec.Command("kubectl", "krew", "upgrade").CombinedOutput()
		} else {
			out, err = exec.Command("kubectl", "krew", "upgrade", args[0]).CombinedOutput()
		}
		if err != nil {
			fmt.Printf("Error upgrading plugin: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginUpgradeCmd)
	RootCmd.AddCommand(pluginCmd)
}
