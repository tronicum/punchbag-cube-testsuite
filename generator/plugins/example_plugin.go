//go:build plugin
// +build plugin

package main

import "punchbag-cube-testsuite/generator/internal/generator"

// ExamplePlugin demonstrates a plugin structure for future extensibility.
// Plugins must implement the ResourceGenerator interface and register themselves.
type ExamplePlugin struct{}

func (p ExamplePlugin) Generate(props map[string]interface{}) string {
	// Implement resource generation logic here
	return "// plugin resource block"
}

// Register the plugin (future dynamic loading)
func init() {
	generator.RegisterPlugin("example", ExamplePlugin{})
}
