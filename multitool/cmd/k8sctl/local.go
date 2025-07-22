// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl local: manage local Kubernetes cluster (minikube, kind, k3d, etc.)
package k8sctl

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Manage local Kubernetes cluster (minikube, kind, k3d, etc.)",
}

var localStartCmd = &cobra.Command{
	Use:   "start [provider]",
	Short: "Start a local Kubernetes cluster (minikube|kind|k3d)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := "minikube"
		if len(args) > 0 {
			provider = args[0]
		}
		var startCmd *exec.Cmd
		switch provider {
		case "minikube":
			startCmd = exec.Command("minikube", "start")
		case "kind":
			startCmd = exec.Command("kind", "create", "cluster")
		case "k3d":
			startCmd = exec.Command("k3d", "cluster", "create")
		default:
			fmt.Printf("Unknown local k8s provider: %s\n", provider)
			return
		}
		out, err := startCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error starting %s: %v\n%s\n", provider, err, string(out))
			return
		}
		fmt.Printf("%s started successfully.\n", provider)
	},
}

var localStopCmd = &cobra.Command{
	Use:   "stop [provider]",
	Short: "Stop a local Kubernetes cluster (minikube|kind|k3d)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := "minikube"
		if len(args) > 0 {
			provider = args[0]
		}
		var stopCmd *exec.Cmd
		switch provider {
		case "minikube":
			stopCmd = exec.Command("minikube", "stop")
		case "kind":
			stopCmd = exec.Command("kind", "delete", "cluster")
		case "k3d":
			stopCmd = exec.Command("k3d", "cluster", "delete")
		default:
			fmt.Printf("Unknown local k8s provider: %s\n", provider)
			return
		}
		out, err := stopCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error stopping %s: %v\n%s\n", provider, err, string(out))
			return
		}
		fmt.Printf("%s stopped successfully.\n", provider)
	},
}

var localStatusCmd = &cobra.Command{
	Use:   "status [provider]",
	Short: "Show status of local Kubernetes cluster (minikube|kind|k3d)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := "minikube"
		if len(args) > 0 {
			provider = args[0]
		}
		var statusCmd *exec.Cmd
		switch provider {
		case "minikube":
			statusCmd = exec.Command("minikube", "status")
		case "kind":
			statusCmd = exec.Command("kind", "get", "clusters")
		case "k3d":
			statusCmd = exec.Command("k3d", "cluster", "list")
		default:
			fmt.Printf("Unknown local k8s provider: %s\n", provider)
			return
		}
		out, err := statusCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error getting status for %s: %v\n%s\n", provider, err, string(out))
			return
		}
		fmt.Print(string(out))
	},
}

func init() {
	localCmd.AddCommand(localStartCmd)
	localCmd.AddCommand(localStopCmd)
	localCmd.AddCommand(localStatusCmd)
	RootCmd.AddCommand(localCmd)
}
