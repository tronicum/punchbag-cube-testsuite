// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// plugin_test.go: tests for k8sctl plugin subcommands
package k8sctl

import (
	"os/exec"
	"strings"
	"testing"
)

func TestKrewPluginList(t *testing.T) {
	cmd := exec.Command("kubectl", "krew", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("kubectl krew list failed: %v\n%s", err, string(out))
	}
	output := string(out)
	// krew list now outputs only plugin names, one per line. Check for at least one known plugin (krew itself)
	if !strings.Contains(output, "krew") {
		t.Errorf("Expected at least 'krew' plugin in output, got: %s", output)
	}
}

func TestKrewPluginInstallAndRemove(t *testing.T) {
	// This is a dry-run style test: do not actually install/remove plugins in CI
	// Instead, check that the install/remove commands are available
	cmd := exec.Command("kubectl", "krew", "help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("kubectl krew help failed: %v\n%s", err, string(out))
	}
	output := string(out)
	if !strings.Contains(output, "install") || !strings.Contains(output, "uninstall") {
		t.Errorf("Expected install/uninstall in krew help output, got: %s", output)
	}
}
