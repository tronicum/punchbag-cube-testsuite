// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl: kubectl-like subcommand for multitool
package k8sctl

import (
	"github.com/spf13/cobra"
)

// RootCmd is the root command for k8sctl (kubectl-like operations)
var RootCmd = &cobra.Command{
	Use:   "k8sctl",
	Short: "kubectl-like operations (get, apply, exec, logs, etc.) via multitool abstraction",
}

func init() {
	AddModeFlag(RootCmd)
}
