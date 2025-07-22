// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl version: kubectl-like version command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version [flags]",
	Short: "Print the client and server version information (kubectl-like)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		kubectlArgs := []string{"version", "--short"}
		if kubeconfig != "" {
			kubectlArgs = append([]string{"--kubeconfig", kubeconfig}, kubectlArgs...)
		}
		cmdName := "kubectl"
		out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl version: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	versionCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
	versionCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	RootCmd.AddCommand(versionCmd)
}
