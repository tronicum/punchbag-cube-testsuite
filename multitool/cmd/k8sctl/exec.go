// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl exec: kubectl-like exec command for multitool
package k8sctl

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [pod] -- [command]",
	Short: "Execute a command in a container (kubectl-like)",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// provider flag is present for future use; defaults to vanilla kubectl unless provider-specific logic is implemented
		// provider, _ := cmd.Flags().GetString("provider")
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		pod := args[0]
		sepIndex := 1
		for i, arg := range args {
			if arg == "--" {
				sepIndex = i
				break
			}
		}
		kubectlArgs := []string{"exec", pod}
		if sepIndex < len(args)-1 {
			kubectlArgs = append(kubectlArgs, args[1:sepIndex+1]...)
			// Only allow alphanumeric and safe shell chars in command
			cmdArgs := args[sepIndex+1:]
			for _, c := range cmdArgs {
				if strings.ContainsAny(c, "|;&><`$") {
					fmt.Println("Unsafe character detected in exec command. Aborting.")
					return
				}
			}
			kubectlArgs = append(kubectlArgs, cmdArgs...)
		}
		if kubeconfig != "" {
			kubectlArgs = append([]string{"--kubeconfig", kubeconfig}, kubectlArgs...)
		}
		cmdName := "kubectl"
		out, err := exec.Command(cmdName, kubectlArgs...).CombinedOutput()
		if err != nil {
			fmt.Printf("Error running kubectl exec: %v\n%s\n", err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	execCmd.Flags().String("provider", "", "Cloud provider (hetzner|azure|...)")
	execCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")
	// No need to add --mode flag here; inherited from RootCmd
	RootCmd.AddCommand(execCmd)
}
