package config

import (
	"fmt"

	"github.com/a7d-corp/i3-config-generator-go/monitor"
)

// Config represents the complete configuration structure
type Config struct {
	I3                  I3Config                `yaml:"i3"`
	UseDetectedMonitors bool                    `yaml:"use_detected_monitors"`
	MonitorDetection    monitor.MonitorConfig   `yaml:"monitor_detection"`
	Layouts             map[string]LayoutConfig `yaml:"layouts"`
	ApplicationBindings map[string]string       `yaml:"application_bindings"`
	StartupPrograms     []string                `yaml:"startup_programs"`
	WindowOverrides     []string                `yaml:"window_overrides"`
	Colors              ColorConfig             `yaml:"colors"`
}

// I3Config holds basic i3 window manager settings
type I3Config struct {
	ModKey  string `yaml:"mod_key"`
	BarFont string `yaml:"bar_font"`
}

// LayoutConfig represents a screen layout configuration (two_mon, one_mon, no_mon)
type LayoutConfig struct {
	GapsInner          int               `yaml:"gaps_inner"`
	GapsOuter          int               `yaml:"gaps_outer"`
	MoveWorkspace      map[string]string `yaml:"move_workspace"`
	WorkspaceToDisplay map[string]string `yaml:"workspace_to_display"`
}

// ColorConfig holds the color scheme configuration
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
	// Basic validation - can be expanded later
	if c.I3.ModKey == "" {
		return fmt.Errorf("i3.mod_key is required")
	}

	if c.UseDetectedMonitors && c.MonitorDetection.DetectionCommand == "" {
		return fmt.Errorf("monitor_detection.detection_command is required when use_detected_monitors is true")
	}

	return nil
}

// GetLayout returns the layout configuration for the given layout name
func (c *Config) GetLayout(layoutName string) (*LayoutConfig, error) {
	layout, exists := c.Layouts[layoutName]
	if !exists {
		return nil, fmt.Errorf("layout '%s' not found in configuration", layoutName)
	}
	return &layout, nil
}
