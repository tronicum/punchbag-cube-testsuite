// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl describe: kubectl-like describe command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe [resource] [name] [flags]",
	Short: "Show details of a specific resource or group of resources (kubectl-like)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		output, _ := cmd.Flags().GetString("output")
		kubectlArgs := []string{"describe"}
		kubectlArgs = append(kubectlArgs, args...)
		if output != "" {
			kubectlArgs = append(kubectlArgs, "-o", output)
		}
		if kubeconfig != "" {
			kubectlArgs = append([]string{"--kubeconfig", kubeconfig}, kubectlArgs...)
		}
		cmdName := "kubectl"
		out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl describe: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	describeCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
	describeCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	describeCmd.Flags().StringP("output", "o", "", "Output format: json|yaml|wide|name|custom-columns|go-template|jsonpath|... (if supported)")
	// No need to add --mode flag here; inherited from RootCmd
	RootCmd.AddCommand(describeCmd)
}
