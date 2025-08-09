package internal

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// ServerConfig holds all server configuration, including storage backends and dummy buckets
// Example YAML:
// storage:
//
//	dummy_buckets:
//	  aws:
//	    - name: test-bucket
//	      region: us-west-2
//	  hetzner:
//	    - name: test-bucket
//	      region: fsn1
//
// ... other config fields ...
type ServerConfig struct {
	Storage struct {
		DummyBuckets map[string][]struct {
			Name   string `yaml:"name"`
			Region string `yaml:"region"`
		} `yaml:"dummy_buckets"`
	} `yaml:"storage"`
	// Add other config fields as needed
	FastSimulate bool `yaml:"fast_simulate"`
	Debug        bool `yaml:"debug"`
}

// LoadServerConfig loads config from conf/config.yaml or path in CUBE_SERVER_CONFIG
func LoadServerConfig() (*ServerConfig, error) {
	path := os.Getenv("CUBE_SERVER_CONFIG")
	if path == "" {
		path = "conf/config.yaml"
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg ServerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// ENV override for fast simulate
	if os.Getenv("FAST_SIMULATE") == "1" {
		cfg.FastSimulate = true
	}
	if os.Getenv("CUBE_SERVER_DEBUG") == "1" {
		cfg.Debug = true
	}
	return &cfg, nil
}
