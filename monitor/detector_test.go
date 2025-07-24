package monitor

import (
	"testing"
)

func TestDetector_padWithDummyMonitors(t *testing.T) {
	config := MonitorConfig{
		DummyMonitors: []string{"dummy1", "dummy2"},
		MinMonitors:   3,
	}
	detector := NewDetector(config)

	tests := []struct {
		name     string
		monitors []string
		expected []string
	}{
		{
			name:     "single monitor - should pad with two dummies",
			monitors: []string{"eDP-1"},
			expected: []string{"eDP-1", "dummy1", "dummy2"},
		},
		{
			name:     "two monitors - should pad with one dummy",
			monitors: []string{"eDP-1", "HDMI-1"},
			expected: []string{"eDP-1", "HDMI-1", "dummy1"},
		},
		{
			name:     "three monitors - no padding needed",
			monitors: []string{"eDP-1", "HDMI-1", "DP-1"},
			expected: []string{"eDP-1", "HDMI-1", "DP-1"},
		},
		{
			name:     "four monitors - no padding needed",
			monitors: []string{"eDP-1", "HDMI-1", "DP-1", "DP-2"},
			expected: []string{"eDP-1", "HDMI-1", "DP-1", "DP-2"},
		},
		{
			name:     "no monitors - should pad with two dummies",
			monitors: []string{},
			expected: []string{"dummy1", "dummy2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.padWithDummyMonitors(tt.monitors)

			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("expected[%d] = %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestDetectedMonitors_GetMonitorByRole(t *testing.T) {
	monitors := &DetectedMonitors{
		Primary: "eDP-1",
		Left:    "HDMI-1",
		Right:   "DP-1",
		All:     []string{"eDP-1", "HDMI-1", "DP-1"},
	}

	tests := []struct {
		role     string
		expected string
	}{
		{"primary_display", "eDP-1"},
		{"left_display", "HDMI-1"},
		{"right_display", "DP-1"},
		{"invalid_role", ""},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			result := monitors.GetMonitorByRole(tt.role)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestDetector_DetectMonitors_Assignment(t *testing.T) {
	// Test the assignment logic without actually running xrandr
	config := MonitorConfig{
		DetectionCommand: "echo 'eDP-1\nHDMI-1'", // Mock command that returns two monitors
		DummyMonitors:    []string{"dummy1", "dummy2"},
		MinMonitors:      3,
	}
	detector := NewDetector(config)

	monitors, err := detector.DetectMonitors()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have eDP-1 as primary, HDMI-1 as left, dummy1 as right
	expected := &DetectedMonitors{
		Primary: "eDP-1",
		Left:    "HDMI-1",
		Right:   "dummy1",
		All:     []string{"eDP-1", "HDMI-1", "dummy1"},
	}

	if monitors.Primary != expected.Primary {
		t.Errorf("expected Primary = %s, got %s", expected.Primary, monitors.Primary)
	}
	if monitors.Left != expected.Left {
		t.Errorf("expected Left = %s, got %s", expected.Left, monitors.Left)
	}
	if monitors.Right != expected.Right {
		t.Errorf("expected Right = %s, got %s", expected.Right, monitors.Right)
	}
}
