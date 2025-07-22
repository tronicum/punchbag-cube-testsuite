// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl get: kubectl-like get command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [resource] [flags]",
	Short: "Get Kubernetes resources (kubectl-like)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
			   kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
			   if kubeconfig == "" {
					   kubeconfig = getKubeconfigForMode()
			   }
			   resource := args[0]
			   kubectlArgs := []string{"--kubeconfig", kubeconfig, "get", resource}
			   if len(args) > 1 {
					   kubectlArgs = append(kubectlArgs, args[1:]...)
			   }
			   cmdName := "kubectl"
			   out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
			   if err != nil {
					   fmt.Printf("Error running kubectl: %v\n%s\n", err, string(out))
					   return
			   }
			   fmt.Print(string(out))
	},
}

func init() {
getCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
getCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
// No need to add --mode flag here; inherited from RootCmd
RootCmd.AddCommand(getCmd)
}
