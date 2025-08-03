package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type MTConfig struct {
	Provider    string                 `yaml:"provider"`
	Region      string                 `yaml:"region"`
	Credentials map[string]string      `yaml:"credentials"`
	Endpoints   map[string]string      `yaml:"endpoints"`
	Settings    map[string]interface{} `yaml:"settings"`
}

func expandEnvOrPrompt(val string) string {
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)\}`)
	return re.ReplaceAllStringFunc(val, func(match string) string {
		varName := re.FindStringSubmatch(match)[1]
		envVal := os.Getenv(varName)
		if envVal != "" {
			return envVal
		}
		// Try to read from a file named after the var
		if f, err := os.ReadFile(varName); err == nil {
			return strings.TrimSpace(string(f))
		}
		// If running in a test, just return the literal
		if os.Getenv("GO_TEST") == "1" {
			return match
		}
		// Prompt the user
		fmt.Printf("Enter value for %s: ", varName)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text()
		}
		return match
	})
}

func expandConfigVars(cfg *MTConfig) {
	cfg.Provider = expandEnvOrPrompt(cfg.Provider)
	cfg.Region = expandEnvOrPrompt(cfg.Region)
	for k, v := range cfg.Credentials {
		cfg.Credentials[k] = expandEnvOrPrompt(v)
	}
	for k, v := range cfg.Endpoints {
		cfg.Endpoints[k] = expandEnvOrPrompt(v)
	}
	for k, v := range cfg.Settings {
		if s, ok := v.(string); ok {
			cfg.Settings[k] = expandEnvOrPrompt(s)
		}
	}
}

func LoadMTConfig(profile string) (*MTConfig, error) {
	base := "multitool/.mtconfig"
	if env := os.Getenv("MTCONFIG_BASE"); env != "" {
		base = env
	}
	if profile == "" {
		profile = "default"
	}
	cfgPath := filepath.Join(base, profile, "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	var cfg MTConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	expandConfigVars(&cfg)
	return &cfg, nil
}
