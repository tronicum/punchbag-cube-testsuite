package importpkg

import "fmt"

// ValidateConfig validates a Config struct (extend with real rules)
func ValidateConfig(cfg *Config) error {
	if cfg.ServerURL == "" {
		return fmt.Errorf("server_url is required")
	}
	if cfg.DefaultProvider == "" {
		return fmt.Errorf("default_provider is required")
	}
	if cfg.DefaultRegion == "" {
		return fmt.Errorf("default_region is required")
	}
	return nil
}
