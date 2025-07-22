// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2023-2025 tronicum@user.github.com
//
// plugin_installed_test.go: tests for k8sctl plugin installed command
package k8sctl

import (
	"os/exec"
	"strings"
	"testing"
)

func TestPluginInstalledCommand(t *testing.T) {
	plugin := "krew" // krew should always be present if krew is installed
	cmd := exec.Command("kubectl", "krew", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("kubectl krew list failed (krew not installed?): %v\n%s", err, string(out))
	}
	output := string(out)
	if !strings.Contains(output, plugin) {
		t.Errorf("Expected plugin '%s' to be listed in krew output, got: %s", plugin, output)
	}
}
