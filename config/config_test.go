package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader_findConfigFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	loader := NewLoader(tempDir)

	// Test: No config file exists
	_, err = loader.findConfigFile()
	if err == nil {
		t.Error("Expected error when no config file exists")
	}

	// Test: Create .yml file
	ymlPath := filepath.Join(tempDir, ConfigFileYML)
	if err := os.WriteFile(ymlPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create yml file: %v", err)
	}

	foundPath, err := loader.findConfigFile()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Should find .yml file if only .yml exists
	expectedPath := ymlPath
	if foundPath != expectedPath {
		t.Errorf("Expected %s, got %s", expectedPath, foundPath)
	}

	// Test: Create .yaml file (should be preferred over .yml)
	yamlPath := filepath.Join(tempDir, ConfigFileYAML)
	if err := os.WriteFile(yamlPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create yaml file: %v", err)
	}

	foundPath, err = loader.findConfigFile()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Should prefer .yaml over .yml
	expectedPath = yamlPath
	if foundPath != expectedPath {
		t.Errorf("Expected %s, got %s", expectedPath, foundPath)
	}
}

func TestLoader_LoadFromFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	loader := NewLoader(tempDir)

	// Test: Valid config file
	validConfig := `
i3:
  mod_key: "Mod4"
  bar_font: "pango:SFNS Display 7, FontAwesome 7"

use_detected_monitors: true

monitor_detection:
  use_native: false
  detection_command: "xrandr | awk '/ connected/' | awk '{print $1}'"
  dummy_monitors:
    - "dummy1"
    - "dummy2"
  min_monitors: 3

layouts:
  two_mon:
    gaps_inner: 20
    gaps_outer: 0
    move_workspace:
      "Ctrl+Shift+1": "left_display"
    workspace_to_display:
      "1": "left_display"

application_bindings:
  "[class=\"^Firefox$\"]": "1"

startup_programs:
  - "/usr/bin/numlockx on"

window_overrides:
  - "[class=\"^[Vv]irtual[Bb]ox*$\"] floating enable"

colors:
  base00: "#1B2B34"
  base01: "#343D46"
`

	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config, err := loader.LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify basic config values
	if config.I3.ModKey != "Mod4" {
		t.Errorf("Expected ModKey 'Mod4', got '%s'", config.I3.ModKey)
	}

	if !config.UseDetectedMonitors {
		t.Error("Expected UseDetectedMonitors to be true")
	}

	if config.MonitorDetection.MinMonitors != 3 {
		t.Errorf("Expected MinMonitors 3, got %d", config.MonitorDetection.MinMonitors)
	}

	// Test layout configuration
	layout, exists := config.Layouts["two_mon"]
	if !exists {
		t.Error("Expected 'two_mon' layout to exist")
	}

	if layout.GapsInner != 20 {
		t.Errorf("Expected GapsInner 20, got %d", layout.GapsInner)
	}

	// Test: Invalid YAML
	invalidPath := filepath.Join(tempDir, "invalid.yaml")
	if err := os.WriteFile(invalidPath, []byte("invalid: yaml: ["), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	_, err = loader.LoadFromFile(invalidPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with native detection",
			config: Config{
				I3:                  I3Config{ModKey: "Mod4"},
				UseDetectedMonitors: true,
				MonitorDetection: MonitorConfig{
					UseNative: true,
					Display:   ":0",
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with shell command detection",
			config: Config{
				I3:                  I3Config{ModKey: "Mod4"},
				UseDetectedMonitors: true,
				MonitorDetection: MonitorConfig{
					UseNative:        false,
					DetectionCommand: "xrandr | awk '/ connected/' | awk '{print $1}'",
				},
			},
			wantErr: false,
		},
		{
			name: "missing mod key",
			config: Config{
				I3: I3Config{ModKey: ""},
			},
			wantErr: true,
		},
		{
			name: "monitor detection enabled but no detection method",
			config: Config{
				I3:                  I3Config{ModKey: "Mod4"},
				UseDetectedMonitors: true,
				MonitorDetection: MonitorConfig{
					UseNative:        false,
					DetectionCommand: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GetLayout(t *testing.T) {
	config := Config{
		Layouts: map[string]LayoutConfig{
			"two_mon": {
				GapsInner: 20,
				GapsOuter: 0,
			},
		},
	}

	// Test: Existing layout
	layout, err := config.GetLayout("two_mon")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if layout.GapsInner != 20 {
		t.Errorf("Expected GapsInner 20, got %d", layout.GapsInner)
	}

	// Test: Non-existing layout
	_, err = config.GetLayout("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existing layout")
	}
}
