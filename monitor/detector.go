package monitor

import (
	"fmt"
	"os/exec"
	"strings"
)

// MonitorConfig holds the configuration for monitor detection
type MonitorConfig struct {
	DetectionCommand string   `yaml:"detection_command"`
	DummyMonitors    []string `yaml:"dummy_monitors"`
	MinMonitors      int      `yaml:"min_monitors"`
}

// DetectedMonitors represents the monitors found on the system
type DetectedMonitors struct {
	Primary string
	Left    string
	Right   string
	All     []string
}

// Detector handles monitor detection operations
type Detector struct {
	config MonitorConfig
}

// NewDetector creates a new monitor detector with the given configuration
func NewDetector(config MonitorConfig) *Detector {
	return &Detector{
		config: config,
	}
}

// DetectMonitors executes the detection command and returns the detected monitors
func (d *Detector) DetectMonitors() (*DetectedMonitors, error) {
	// Execute the detection command
	monitors, err := d.executeDetectionCommand()
	if err != nil {
		return nil, fmt.Errorf("failed to detect monitors: %w", err)
	}

	// Ensure we have the minimum number of monitors by padding with dummies
	paddedMonitors := d.padWithDummyMonitors(monitors)

	// Assign monitors to roles
	result := &DetectedMonitors{
		All: paddedMonitors,
	}

	// Assign primary, left, and right displays based on array indices
	if len(paddedMonitors) >= 1 {
		result.Primary = paddedMonitors[0]
	}
	if len(paddedMonitors) >= 2 {
		result.Left = paddedMonitors[1]
	}
	if len(paddedMonitors) >= 3 {
		result.Right = paddedMonitors[2]
	}

	return result, nil
}

// executeDetectionCommand runs the configured detection command and parses the output
func (d *Detector) executeDetectionCommand() ([]string, error) {
	// Execute the command using shell
	cmd := exec.Command("sh", "-c", d.config.DetectionCommand)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	// Parse the output - split by newlines and filter out empty strings
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var monitors []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			monitors = append(monitors, line)
		}
	}

	return monitors, nil
}

// padWithDummyMonitors ensures we have at least the minimum number of monitors
// by adding dummy monitors as needed
func (d *Detector) padWithDummyMonitors(monitors []string) []string {
	result := make([]string, len(monitors))
	copy(result, monitors)

	// Add dummy monitors if we don't have enough
	dummyIndex := 0
	for len(result) < d.config.MinMonitors && dummyIndex < len(d.config.DummyMonitors) {
		result = append(result, d.config.DummyMonitors[dummyIndex])
		dummyIndex++
	}

	return result
}

// GetMonitorByRole returns the monitor name for a given role (primary_display, left_display, right_display)
func (dm *DetectedMonitors) GetMonitorByRole(role string) string {
	switch role {
	case "primary_display":
		return dm.Primary
	case "left_display":
		return dm.Left
	case "right_display":
		return dm.Right
	default:
		return ""
	}
}

// String returns a string representation of the detected monitors
func (dm *DetectedMonitors) String() string {
	return fmt.Sprintf("Primary: %s, Left: %s, Right: %s, All: %v",
		dm.Primary, dm.Left, dm.Right, dm.All)
}
