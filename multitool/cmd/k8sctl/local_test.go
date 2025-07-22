// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// local_test.go: tests for k8sctl local subcommands
package k8sctl

import (
	"os/exec"
	"testing"
)

func TestLocalStatusMinikube(t *testing.T) {
	cmd := exec.Command("minikube", "status")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skip("minikube not installed or not running")
	}
}

func TestLocalStatusKind(t *testing.T) {
	cmd := exec.Command("kind", "get", "clusters")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skip("kind not installed or not running")
	}
}

func TestLocalStatusK3d(t *testing.T) {
	cmd := exec.Command("k3d", "cluster", "list")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Skip("k3d not installed or not running")
	}
}
