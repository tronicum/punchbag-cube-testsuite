package main

import "fmt"

// safeString returns a string value from a map or a default
func safeString(m map[string]interface{}, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

// safeInt returns an int value from a map or a default
func safeInt(m map[string]interface{}, key string, def int) int {
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
func safeBool(m map[string]interface{}, key string, def bool) bool {
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
func logPrefix(context string) string {
	return fmt.Sprintf("[werfty][%s]", context)
}
