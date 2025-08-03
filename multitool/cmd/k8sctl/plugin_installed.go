// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl plugin installed: check if a kubectl plugin is installed via krew
package k8sctl

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var pluginInstalledCmd = &cobra.Command{
	Use:   "installed [plugin]",
	Short: "Check if a kubectl plugin is installed via krew",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plugin := args[0]
		out, err := exec.Command("kubectl", "krew", "list").CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl krew list: %v\n%s\n", err, string(out))
			return
		}
		output := string(out)
		if strings.Contains(output, "\n"+plugin+"\n") || strings.HasPrefix(output, plugin+"\n") || strings.Contains(output, "\n"+plugin+" ") {
			fmt.Printf("Plugin '%s' is installed.\n", plugin)
		} else {
			fmt.Printf("Plugin '%s' is NOT installed.\n", plugin)
		}
	},
}

func init() {
	pluginCmd.AddCommand(pluginInstalledCmd)
}
