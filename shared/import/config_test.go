package importpkg

import (
	"strings"
	"testing"
)

func TestLoadConfigJSONAndValidate(t *testing.T) {
	jsonData := `{
		"server_url": "http://localhost:8080",
		"default_provider": "azure",
		"default_region": "westeurope"
	}`
	r := strings.NewReader(jsonData)
	cfg, err := LoadConfigJSON(r)
	if err != nil {
		t.Fatalf("LoadConfigJSON failed: %v", err)
	}
	if err := ValidateConfig(cfg); err != nil {
		t.Fatalf("ValidateConfig failed: %v", err)
	}
}

func TestLoadConfigYAMLAndValidate(t *testing.T) {
	yamlData := `
server_url: http://localhost:8080
default_provider: azure
default_region: westeurope
`
	r := strings.NewReader(yamlData)
	cfg, err := LoadConfigYAML(r)
	if err != nil {
		t.Fatalf("LoadConfigYAML failed: %v", err)
	}
	if err := ValidateConfig(cfg); err != nil {
		t.Fatalf("ValidateConfig failed: %v", err)
	}
}

func TestValidateConfigFails(t *testing.T) {
	cfg := &Config{}
	if err := ValidateConfig(cfg); err == nil {
		t.Fatal("expected validation to fail for empty config")
	}
}
