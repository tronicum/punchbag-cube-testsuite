// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl autoscale: kubectl-like autoscale command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var autoscaleCmd = &cobra.Command{
	Use:   "autoscale [resource] [flags]",
	Short: "Auto-scale a deployment, replica set, stateful set, or replication controller (kubectl-like)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		kubectlArgs := []string{"autoscale"}
		kubectlArgs = append(kubectlArgs, args...)
		if kubeconfig != "" {
			kubectlArgs = append([]string{"--kubeconfig", kubeconfig}, kubectlArgs...)
		}
		cmdName := "kubectl"
		out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl autoscale: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	autoscaleCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
	autoscaleCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	RootCmd.AddCommand(autoscaleCmd)
}
