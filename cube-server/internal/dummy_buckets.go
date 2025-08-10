package internal

import (
	"encoding/json"
	"github.com/tronicum/punchbag-cube-testsuite/shared/simulation"
	"os"
	"path/filepath"
	"strings"
)

// InjectDummyBucketsIfNeeded injects dummy buckets for each provider if no buckets exist yet
func InjectDummyBucketsIfNeeded(sim *simulation.SimulationService, dummy map[string][]struct{ Name, Region string }) {
	for provider, buckets := range dummy {
		if len(sim.BucketStore().List(provider)) == 0 {
			for _, b := range buckets {
				sim.BucketStore().Create(provider, b.Name, b.Region)
			}
		}
	}
}

// DummyBucketsConfig holds configuration for dummy/test buckets
// Can be loaded from file or ENV
// Example JSON: { "aws": [ { "name": "test-bucket", "region": "us-west-2" } ] }
type DummyBucketsConfig map[string][]DummyBucket

type DummyBucket struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

// LoadDummyBucketsConfig loads dummy bucket config from file, ENV, or detects simulation mode by path
func LoadDummyBucketsConfig() DummyBucketsConfig {
	// 1. Check ENV var
	if env := os.Getenv("CUBE_SERVER_DUMMY_BUCKETS"); env != "" {
		var cfg DummyBucketsConfig
		_ = json.Unmarshal([]byte(env), &cfg)
		return cfg
	}
	// 2. Check config file (default: ./conf/dummy_buckets.json)
	confPath := os.Getenv("CUBE_SERVER_DUMMY_BUCKETS_FILE")
	if confPath == "" {
		confPath = filepath.Join("conf", "dummy_buckets.json")
	}
	if data, err := os.ReadFile(confPath); err == nil {
		var cfg DummyBucketsConfig
		if err := json.Unmarshal(data, &cfg); err == nil {
			return cfg
		}
	}
	// 3. Detect simulation mode by path (if running in a known test/sim dir)
	if strings.Contains(os.Args[0], "sim") || strings.Contains(os.Args[0], "test") {
		return DummyBucketsConfig{
			"aws":     {{Name: "auto-sim-bucket", Region: "us-west-2"}},
			"hetzner": {{Name: "auto-sim-bucket", Region: "fsn1"}},
		}
	}
	return DummyBucketsConfig{}
}

// InjectDummyBuckets injects dummy buckets into the simulation service for all configured providers
func InjectDummyBuckets(sim *simulation.SimulationService, cfg DummyBucketsConfig) {
	for provider, buckets := range cfg {
		for _, b := range buckets {
			sim.BucketStore().Create(provider, b.Name, b.Region)
		}
	}
}
