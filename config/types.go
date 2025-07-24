package config

import (
	"fmt"

	"github.com/a7d-corp/i3-config-generator-go/monitor"
)

type Config struct {
	I3                  I3Config                `yaml:"i3"`
	UseDetectedMonitors bool                    `yaml:"use_detected_monitors"`
	MonitorDetection    MonitorConfig           `yaml:"monitor_detection"`
	Layouts             map[string]LayoutConfig `yaml:"layouts"`
	ApplicationBindings map[string]string       `yaml:"application_bindings"`
	StartupPrograms     []string                `yaml:"startup_programs"`
	WindowOverrides     []string                `yaml:"window_overrides"`
	Colors              ColorConfig             `yaml:"colors"`
}

type I3Config struct {
	ModKey  string `yaml:"mod_key"`
	BarFont string `yaml:"bar_font"`
}

type MonitorConfig struct {
	// Native X11 detection settings (preferred)
	UseNative bool   `yaml:"use_native"`
	Display   string `yaml:"display"`

	// Legacy shell command detection (deprecated but supported)
	DetectionCommand string `yaml:"detection_command"`

	// Common settings for both approaches
	DummyMonitors []string `yaml:"dummy_monitors"`
	MinMonitors   int      `yaml:"min_monitors"`
}

type LayoutConfig struct {
	GapsInner          int               `yaml:"gaps_inner"`
	GapsOuter          int               `yaml:"gaps_outer"`
	MoveWorkspace      map[string]string `yaml:"move_workspace"`
	WorkspaceToDisplay map[string]string `yaml:"workspace_to_display"`
}

type ColorConfig struct {
	Base00 string `yaml:"base00"`
	Base01 string `yaml:"base01"`
	Base02 string `yaml:"base02"`
	Base03 string `yaml:"base03"`
	Base04 string `yaml:"base04"`
	Base05 string `yaml:"base05"`
	Base06 string `yaml:"base06"`
	Base07 string `yaml:"base07"`
	Base08 string `yaml:"base08"`
	Base09 string `yaml:"base09"`
	Base0A string `yaml:"base0A"`
	Base0B string `yaml:"base0B"`
	Base0C string `yaml:"base0C"`
	Base0D string `yaml:"base0D"`
	Base0E string `yaml:"base0E"`
	Base0F string `yaml:"base0F"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.I3.ModKey == "" {
		return fmt.Errorf("i3.mod_key is required")
	}

	if c.UseDetectedMonitors {
		// Check if either native detection is enabled or command detection is configured
		if !c.MonitorDetection.UseNative && c.MonitorDetection.DetectionCommand == "" {
			return fmt.Errorf("monitor detection is enabled but neither use_native nor detection_command is configured")
		}
	}

	return nil
}

// GetLayout returns the layout configuration for the given name
func (c *Config) GetLayout(layoutName string) (*LayoutConfig, error) {
	layout, exists := c.Layouts[layoutName]
	if !exists {
		return nil, fmt.Errorf("layout '%s' not found", layoutName)
	}
	return &layout, nil
}

// CreateDetector creates an appropriate monitor detector based on configuration
func (c *Config) CreateDetector() (monitor.MonitorDetector, error) {
	if !c.UseDetectedMonitors {
		return nil, fmt.Errorf("monitor detection is disabled")
	}

	// Prefer native detection
	if c.MonitorDetection.UseNative {
		return monitor.NewNativeDetector(
			c.MonitorDetection.Display,
			c.MonitorDetection.DummyMonitors,
			c.MonitorDetection.MinMonitors,
		), nil
	}

	// Fall back to shell command detection if configured
	if c.MonitorDetection.DetectionCommand != "" {
		// Convert to the old MonitorConfig format for backward compatibility
		oldConfig := monitor.MonitorConfig{
			DetectionCommand: c.MonitorDetection.DetectionCommand,
			DummyMonitors:    c.MonitorDetection.DummyMonitors,
			MinMonitors:      c.MonitorDetection.MinMonitors,
		}
		return monitor.NewDetector(oldConfig), nil
	}

	return nil, fmt.Errorf("no valid monitor detection method configured")
}
