// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// k8sctl krew: manage krew plugin manager for kubectl plugins
package k8sctl

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var krewCmd = &cobra.Command{
	Use:   "krew",
	Short: "Manage krew plugin manager (install, upgrade, etc.)",
}

var krewInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install krew plugin manager (requires curl, tar, and unzip)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// This is a simplified version of the official install script
		fmt.Println("Installing krew plugin manager...")
		cmdStr := `set -e; cd "$(mktemp -d)" && \
  OS="$(uname | tr '[:upper:]' '[:lower:]')" && \
  ARCH="$(uname -m | sed 's/x86_64/amd64/;s/arm.*$/arm/;s/aarch64$/arm64/")" && \
  KREW="krew-${OS}_${ARCH}" && \
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" && \
  tar zxvf "${KREW}.tar.gz" && \
  ./${KREW} install krew`
		installCmd := exec.Command("sh", "-c", cmdStr)
		out, err := installCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error installing krew: %v\n%s\n", err, string(out))
			return
		}
		fmt.Println("krew installed successfully.")
	},
}

func init() {
	krewCmd.AddCommand(krewInstallCmd)
	RootCmd.AddCommand(krewCmd)
}
