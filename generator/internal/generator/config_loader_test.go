package generator

import (
	"os"
	"testing"
)

func TestLoadConfigFromFile_JSON(t *testing.T) {
	jsonData := `{"name": "test-resource", "properties": {"foo": "bar"}}`
	file := "test_config.json"
	os.WriteFile(file, []byte(jsonData), 0644)
	defer os.Remove(file)
	cfg, err := LoadConfigFromFile(file)
	if err != nil {
		t.Fatalf("failed to load JSON config: %v", err)
	}
	if cfg["name"] != "test-resource" {
		t.Errorf("expected name 'test-resource', got %v", cfg["name"])
	}
}

func TestLoadConfigFromFile_YAML(t *testing.T) {
	yamlData := `name: test-resource
properties:
  foo: bar
`
	file := "test_config.yaml"
	os.WriteFile(file, []byte(yamlData), 0644)
	defer os.Remove(file)
	cfg, err := LoadConfigFromFile(file)
	if err != nil {
		t.Fatalf("failed to load YAML config: %v", err)
	}
	if cfg["name"] != "test-resource" {
		t.Errorf("expected name 'test-resource', got %v", cfg["name"])
	}
}
