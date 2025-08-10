// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl label: kubectl-like label command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:   "label [TYPE/NAME] [KEY=VALUE ...] [flags]",
	Short: "Update the labels on resources (kubectl-like)",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		kubectlArgs := []string{"label"}
		kubectlArgs = append(kubectlArgs, args...)
		if kubeconfig != "" {
			kubectlArgs = append([]string{"--kubeconfig", kubeconfig}, kubectlArgs...)
		}
		cmdName := "kubectl"
		out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl label: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	labelCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
	labelCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	RootCmd.AddCommand(labelCmd)
}
