package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Default configuration directory relative to user's home
	DefaultConfigDir = ".config/i3-config-generator"
	// Configuration file names to search for (in order of preference)
	ConfigFileYAML = "config.yaml"
	ConfigFileYML  = "config.yml"
)

// Loader handles configuration file loading operations
type Loader struct {
	configDir string
}

// NewLoader creates a new configuration loader
// If configDir is empty, uses the default location ($HOME/.config/i3-config-generator)
func NewLoader(configDir string) *Loader {
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory if home directory is not available
			configDir = DefaultConfigDir
		} else {
			configDir = filepath.Join(homeDir, DefaultConfigDir)
		}
	}

	return &Loader{
		configDir: configDir,
	}
}

// Load attempts to load the configuration file from the configured directory
// It tries both .yaml and .yml extensions in that order
func (l *Loader) Load() (*Config, error) {
	configPath, err := l.findConfigFile()
	if err != nil {
		return nil, err
	}

	return l.LoadFromFile(configPath)
}

// LoadFromFile loads configuration from a specific file path
func (l *Loader) LoadFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config file %s: %w", filePath, err)
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// findConfigFile searches for the configuration file in the configured directory
// Returns the path to the first file found (checking .yaml first, then .yml)
func (l *Loader) findConfigFile() (string, error) {
	// List of filenames to try in order of preference
	candidates := []string{ConfigFileYAML, ConfigFileYML}

	for _, filename := range candidates {
		configPath := filepath.Join(l.configDir, filename)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	// If no config file found, return an informative error
	return "", fmt.Errorf("no configuration file found in %s (tried: %s, %s)",
		l.configDir, ConfigFileYAML, ConfigFileYML)
}

// GetConfigDir returns the configuration directory being used
func (l *Loader) GetConfigDir() string {
	return l.configDir
}

// EnsureConfigDir creates the configuration directory if it doesn't exist
func (l *Loader) EnsureConfigDir() error {
	if err := os.MkdirAll(l.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", l.configDir, err)
	}
	return nil
}

// GetConfigPath returns the full path to the preferred config file location
func (l *Loader) GetConfigPath() string {
	return filepath.Join(l.configDir, ConfigFileYAML)
}
