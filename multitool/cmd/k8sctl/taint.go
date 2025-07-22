// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl taint: kubectl-like taint command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var taintCmd = &cobra.Command{
	Use:   "taint [node] [flags]",
	Short: "Update the taints on one or more nodes (kubectl-like)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		kubectlArgs := []string{"taint"}
		kubectlArgs = append(kubectlArgs, args...)
		if kubeconfig != "" {
			kubectlArgs = append([]string{"--kubeconfig", kubeconfig}, kubectlArgs...)
		}
		cmdName := "kubectl"
		out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl taint: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	taintCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
	taintCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	RootCmd.AddCommand(taintCmd)
}
