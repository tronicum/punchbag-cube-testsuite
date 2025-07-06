package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"gopkg.in/yaml.v3"
)

// LoadConfigFromFile loads a resource definition from a YAML or JSON file
func LoadConfigFromFile(path string) (map[string]interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var data map[string]interface{}
	if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
		err = yaml.Unmarshal(content, &data)
		if err != nil {
			return nil, fmt.Errorf("invalid YAML: %w", err)
		}
	} else {
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
	}
	return data, nil
}
