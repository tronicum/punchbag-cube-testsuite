package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/output"
	"gopkg.in/yaml.v2"
)

// Config represents the multitool configuration
type Config struct {
	ServerURL       string             `json:"server_url" yaml:"server_url"`
	DefaultProvider string             `json:"default_provider" yaml:"default_provider"`
	DefaultRegion   string             `json:"default_region" yaml:"default_region"`
	DefaultOutput   string             `json:"default_output" yaml:"default_output"`
	Profiles        map[string]Profile `json:"profiles" yaml:"profiles"`
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

var (
	configDir           = filepath.Join(os.Getenv("HOME"), ".multitool")
	multitoolConfigFile = filepath.Join(configDir, "config.yaml")
)

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage multitool configuration",
	Long:  "Manage configuration settings, profiles, and defaults for multitool.",
}

// configInitCmd initializes a new configuration
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize multitool configuration",
	Long:  "Initialize a new multitool configuration file with default settings.",
	Run: func(cmd *cobra.Command, args []string) {
		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			output.FormatError(fmt.Errorf("failed to create config directory: %w", err))
			return
		}

		// Check if config already exists
		if _, err := os.Stat(multitoolConfigFile); err == nil {
			overwrite, _ := cmd.Flags().GetBool("force")
			if !overwrite {
				output.FormatWarning("Configuration file already exists. Use --force to overwrite.")
				return
			}
		}

		// Create default configuration
		config := &Config{
			ServerURL:       "http://localhost:8080",
			DefaultProvider: "azure",
			DefaultRegion:   "eastus",
			DefaultOutput:   "table",
			Profiles: map[string]Profile{
				"default": {
					ServerURL: "http://localhost:8080",
					Provider:  "azure",
					Region:    "eastus",
				},
				"development": {
					ServerURL: "http://localhost:8080",
					Provider:  "azure",
					Region:    "eastus",
					Tags: map[string]string{
						"environment": "development",
						"team":        "devops",
					},
				},
				"production": {
					ServerURL: "https://punchbag-prod.example.com",
					Provider:  "azure",
					Region:    "eastus",
					Tags: map[string]string{
						"environment": "production",
						"team":        "devops",
					},
				},
			},
		}

		// Write configuration to file
		data, err := yaml.Marshal(config)
		if err != nil {
			output.FormatError(fmt.Errorf("failed to marshal config: %w", err))
			return
		}

		if err := os.WriteFile(multitoolConfigFile, data, 0644); err != nil {
			output.FormatError(fmt.Errorf("failed to write config file: %w", err))
			return
		}

		output.FormatSuccess(fmt.Sprintf("Configuration initialized at %s", multitoolConfigFile))
	},
}

// configShowCmd shows the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show [profile]",
	Short: "Show current configuration",
	Long:  "Display the current configuration or a specific profile.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			output.FormatError(fmt.Errorf("failed to load config: %w", err))
			return
		}

		var displayData interface{}

		if len(args) > 0 {
			profileName := args[0]
			profile, exists := config.Profiles[profileName]
			if !exists {
				output.FormatError(fmt.Errorf("profile '%s' not found", profileName))
				return
			}
			displayData = profile
		} else {
			displayData = config
		}

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(displayData); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
	},
}

// configSetCmd sets configuration values
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value in the default profile or a specific profile.

Examples:
  multitool config set server_url http://localhost:8080
  multitool config set --profile development provider aws
  multitool config set --profile production region us-west-2`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		profileName, _ := cmd.Flags().GetString("profile")

		config, err := loadConfig()
		if err != nil {
			// Create new config if it doesn't exist
			config = &Config{
				Profiles: make(map[string]Profile),
			}
		}

		if profileName == "" {
			// Set global config
			switch key {
			case "server_url":
				config.ServerURL = value
			case "default_provider":
				config.DefaultProvider = value
			case "default_region":
				config.DefaultRegion = value
			case "default_output":
				config.DefaultOutput = value
			default:
				output.FormatError(fmt.Errorf("unknown config key: %s", key))
				return
			}
		} else {
			// Set profile-specific config
			profile, exists := config.Profiles[profileName]
			if !exists {
				profile = Profile{}
			}

			switch key {
			case "server_url":
				profile.ServerURL = value
			case "provider":
				profile.Provider = value
			case "region":
				profile.Region = value
			case "resource_group":
				profile.ResourceGroup = value
			case "project_id":
				profile.ProjectID = value
			case "location":
				profile.Location = value
			default:
				output.FormatError(fmt.Errorf("unknown profile config key: %s", key))
				return
			}

			config.Profiles[profileName] = profile
		}

		if err := saveConfig(config); err != nil {
			output.FormatError(fmt.Errorf("failed to save config: %w", err))
			return
		}

		if profileName == "" {
			output.FormatSuccess(fmt.Sprintf("Set %s = %s", key, value))
		} else {
			output.FormatSuccess(fmt.Sprintf("Set %s = %s for profile '%s'", key, value, profileName))
		}
	},
}

// configListProfilesCmd lists all profiles
var configListProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List all configuration profiles",
	Long:  "List all available configuration profiles.",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			output.FormatError(fmt.Errorf("failed to load config: %w", err))
			return
		}

		if len(config.Profiles) == 0 {
			output.FormatInfo("No profiles found")
			return
		}

		// Convert profiles to a format suitable for table output
		var profileData []map[string]interface{}
		for name, profile := range config.Profiles {
			profileData = append(profileData, map[string]interface{}{
				"Name":          name,
				"ServerURL":     profile.ServerURL,
				"Provider":      profile.Provider,
				"Region":        profile.Region,
				"ResourceGroup": profile.ResourceGroup,
				"ProjectID":     profile.ProjectID,
			})
		}

		formatter := output.NewFormatter(output.Format(outputFormat))
		if err := formatter.FormatOutput(profileData); err != nil {
			output.FormatError(fmt.Errorf("failed to format output: %w", err))
		}
	},
}

// Helper functions

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(multitoolConfigFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(multitoolConfigFile, data, 0644)
}

// GetProfile loads a specific profile or returns default values
func GetProfile(profileName string) (*Profile, error) {
	config, err := loadConfig()
	if err != nil {
		// Return default profile if config doesn't exist
		return &Profile{
			ServerURL: "http://localhost:8080",
			Provider:  "azure",
			Region:    "eastus",
		}, nil
	}

	if profileName == "" {
		profileName = "default"
	}

	profile, exists := config.Profiles[profileName]
	if !exists {
		return nil, fmt.Errorf("profile '%s' not found", profileName)
	}

	return &profile, nil
}

func init() {
	// Add config subcommands
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListProfilesCmd)

	// Config flags
	configCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	configInitCmd.Flags().Bool("force", false, "Overwrite existing configuration")
	configSetCmd.Flags().String("profile", "", "Profile name to set config for")
}
