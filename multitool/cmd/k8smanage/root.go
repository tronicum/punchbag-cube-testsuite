// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8s-manage: cluster lifecycle management for multitool
package k8smanage

import (
	"github.com/spf13/cobra"
)

// RootCmd is the root command for k8s-manage (cluster lifecycle management)
var RootCmd = &cobra.Command{
	Use:   "k8s-manage",
	Short: "Cluster lifecycle management (create, delete, scale, upgrade, etc.) for all providers",
}

// Add subcommands here (e.g., create cluster, delete cluster, scale, upgrade, etc.)
