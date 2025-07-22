package cmd

import (
	"github.com/spf13/cobra"
	k8sctl "github.com/tronicum/punchbag-cube-testsuite/multitool/cmd/k8sctl"
)

// k8sctlCmd is the top-level kubectl-like command
var k8sctlCmd = &cobra.Command{
	Use:   "k8sctl",
	Short: "kubectl-like operations (get, apply, exec, logs, etc.) via multitool abstraction",
}

func init() {
	k8sctlCmd.AddCommand(k8sctl.RootCmd.Commands()...)
	rootCmd.AddCommand(k8sctlCmd)
}
