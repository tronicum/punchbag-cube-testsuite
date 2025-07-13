package importpkg

import (
  "encoding/json"
  "io"
  "gopkg.in/yaml.v3"
)
// LoadConfigYAML loads configuration from a YAML reader
func LoadConfigYAML(r io.Reader) (*Config, error) {
  var cfg Config
  dec := yaml.NewDecoder(r)
  if err := dec.Decode(&cfg); err != nil {
	return nil, err
  }
  return &cfg, nil
}

// Profile represents a configuration profile for different environments
type Profile struct {
	ServerURL     string            `json:"server_url" yaml:"server_url"`
	Provider      string            `json:"provider" yaml:"provider"`
	Region        string            `json:"region" yaml:"region"`
	ResourceGroup string            `json:"resource_group,omitempty" yaml:"resource_group,omitempty"`
	ProjectID     string            `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	Location      string            `json:"location,omitempty" yaml:"location,omitempty"`
	Tags          map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// TestConfig represents test-specific configuration
type TestConfig struct {
	Duration     string            `json:"duration,omitempty" yaml:"duration,omitempty"`
	Concurrency  int               `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	RequestRate  int               `json:"request_rate,omitempty" yaml:"request_rate,omitempty"`
	TargetURL    string            `json:"target_url,omitempty" yaml:"target_url,omitempty"`
	Method       string            `json:"method,omitempty" yaml:"method,omitempty"`
	ExpectedCode int               `json:"expected_code,omitempty" yaml:"expected_code,omitempty"`
	Headers      map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Body         string            `json:"body,omitempty" yaml:"body,omitempty"`
}

// Config represents the unified configuration structure
type Config struct {
	ServerURL       string             `json:"server_url" yaml:"server_url"`
	DefaultProvider string             `json:"default_provider" yaml:"default_provider"`
	DefaultRegion   string             `json:"default_region" yaml:"default_region"`
	DefaultOutput   string             `json:"default_output" yaml:"default_output"`
	Profiles        map[string]Profile `json:"profiles,omitempty" yaml:"profiles,omitempty"`
	Test            *TestConfig        `json:"test,omitempty" yaml:"test,omitempty"`
}

// LoadConfigJSON loads configuration from a JSON reader
func LoadConfigJSON(r io.Reader) (*Config, error) {
	var cfg Config
	dec := json.NewDecoder(r)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
