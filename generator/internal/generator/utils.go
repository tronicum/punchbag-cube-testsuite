package generator

import "fmt"

// safeString returns a string value from a map or a default
func SafeString(m map[string]interface{}, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

// safeInt returns an int value from a map or a default
func SafeInt(m map[string]interface{}, key string, def int) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int:
			return val
		case float64:
			return int(val)
		}
	}
	return def
}

// safeBool returns a bool value from a map or a default
func SafeBool(m map[string]interface{}, key string, def bool) bool {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case bool:
			return val
		case string:
			return val == "true" || val == "1"
		}
	}
	return def
}

// logPrefix returns a formatted prefix for log messages
func LogPrefix(context string) string {
	return fmt.Sprintf("[werfty][%s]", context)
}

// ResourceGenerator is the interface for all resource plugins
// (core and external)
type ResourceGenerator interface {
	Generate(props map[string]interface{}) string
}

// pluginRegistry holds all registered plugins by name
var pluginRegistry = map[string]ResourceGenerator{}

// RegisterPlugin registers a plugin by name
func RegisterPlugin(name string, plugin ResourceGenerator) {
	pluginRegistry[name] = plugin
}

// GetPlugin returns a registered plugin by name (or nil)
func GetPlugin(name string) ResourceGenerator {
	return pluginRegistry[name]
}

// Example usage of plugin system in orchestration:
// If a plugin is registered for a resourceType, use it for generation.
func GenerateTerraformWithPlugins(resourceType string, props map[string]interface{}) (string, error) {
	plugin := GetPlugin(resourceType)
	if plugin != nil {
		return plugin.Generate(props), nil
	}
	return "", fmt.Errorf("no plugin registered for resourceType: %s", resourceType)
}
