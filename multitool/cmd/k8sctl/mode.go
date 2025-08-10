// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// mode.go: support for --mode=direct|proxy|local in k8sctl
package k8sctl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mode string
var defaultMode string

func getKubeconfigForMode() string {
	switch mode {
	case "local":
		return os.ExpandEnv("$HOME/.kube/config") // local cluster
	case "proxy":
		return os.ExpandEnv("$HOME/.kube/proxy-config") // proxy config (customize as needed)
	case "direct":
		fallthrough
	default:
		return os.ExpandEnv("$HOME/.kube/config") // default remote
	}
}

func AddModeFlag(cmd *cobra.Command) {
	v := viper.New()
	v.SetConfigType("yaml")
	// 1. Project config (./conf/k8sctl.yml)
	v.SetConfigName("k8sctl")
	v.AddConfigPath("./conf")
	v.AddConfigPath("../conf")
	v.AddConfigPath("../../conf")
	// 2. User config ($HOME/.mt/config.yaml)
	userConfig := filepath.Join(os.Getenv("HOME"), ".mt", "config.yaml")
	v.SetDefault("default_mode", "direct")
	v.SetDefault("default_provider", "")
	// 3. ENV overrides (K8SCTL_MODE, K8SCTL_PROVIDER)
	v.AutomaticEnv()
	v.BindEnv("default_mode", "K8SCTL_MODE")
	v.BindEnv("default_provider", "K8SCTL_PROVIDER")
	// Read project config if present
	_ = v.ReadInConfig()
	// Read user config if present (overrides project config)
	if _, err := os.Stat(userConfig); err == nil {
		v.SetConfigFile(userConfig)
		_ = v.MergeInConfig()
	}
	// CLI flag will override all
	defaultMode = v.GetString("default_mode")
	// Register --mode flag (overrides all)
	cmd.PersistentFlags().StringVar(&mode, "mode", defaultMode, "Kubernetes access mode: direct (remote), proxy (via cube proxy), or local (127.0.0.1/minikube/kind/k3d)")
	cobra.OnInitialize(func() {
		fmt.Printf("k8sctl running in '%s' mode\n", mode)
	})
}
