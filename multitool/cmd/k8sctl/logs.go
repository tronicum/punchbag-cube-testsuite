// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl logs: kubectl-like logs command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [pod] [flags]",
	Short: "Fetch logs from a pod (kubectl-like)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		pod := args[0]
		kubectlArgs := []string{"logs", pod}
		if len(args) > 1 {
			kubectlArgs = append(kubectlArgs, args[1:]...)
		}
		if kubeconfig != "" {
			kubectlArgs = append([]string{"--kubeconfig", kubeconfig}, kubectlArgs...)
		}
		cmdName := "kubectl"
		out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl logs: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	logsCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
	logsCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	// No need to add --mode flag here; inherited from RootCmd
	RootCmd.AddCommand(logsCmd)
}
