package monitor

import (
	"testing"
)

func TestNativeDetector_padWithDummyMonitors(t *testing.T) {
	tests := []struct {
		name          string
		monitors      []string
		dummyMonitors []string
		minMonitors   int
		expected      []string
	}{
		{
			name:          "single monitor - should pad with two dummies",
			monitors:      []string{"eDP-1"},
			dummyMonitors: []string{"dummy1", "dummy2"},
			minMonitors:   3,
			expected:      []string{"eDP-1", "dummy1", "dummy2"},
		},
		{
			name:          "two monitors - should pad with one dummy",
			monitors:      []string{"eDP-1", "HDMI-1"},
			dummyMonitors: []string{"dummy1", "dummy2"},
			minMonitors:   3,
			expected:      []string{"eDP-1", "HDMI-1", "dummy1"},
		},
		{
			name:          "three monitors - no padding needed",
			monitors:      []string{"eDP-1", "HDMI-1", "DP-1"},
			dummyMonitors: []string{"dummy1", "dummy2"},
			minMonitors:   3,
			expected:      []string{"eDP-1", "HDMI-1", "DP-1"},
		},
		{
			name:          "four monitors - no padding needed",
			monitors:      []string{"eDP-1", "HDMI-1", "DP-1", "DP-2"},
			dummyMonitors: []string{"dummy1", "dummy2"},
			minMonitors:   3,
			expected:      []string{"eDP-1", "HDMI-1", "DP-1", "DP-2"},
		},
		{
			name:          "no monitors - should pad with two dummies",
			monitors:      []string{},
			dummyMonitors: []string{"dummy1", "dummy2"},
			minMonitors:   3,
			expected:      []string{"dummy1", "dummy2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd := NewNativeDetector(":0", tt.dummyMonitors, tt.minMonitors)
			result := nd.padWithDummyMonitors(tt.monitors)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d monitors, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected monitor[%d] = %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestNewNativeDetector(t *testing.T) {
	tests := []struct {
		name                string
		display             string
		dummyMonitors       []string
		minMonitors         int
		expectedDisplay     string
		expectedMinMonitors int
	}{
		{
			name:                "default values",
			display:             "",
			dummyMonitors:       []string{"dummy1", "dummy2"},
			minMonitors:         0,
			expectedDisplay:     ":0",
			expectedMinMonitors: 3,
		},
		{
			name:                "custom values",
			display:             ":1",
			dummyMonitors:       []string{"dummy1", "dummy2"},
			minMonitors:         2,
			expectedDisplay:     ":1",
			expectedMinMonitors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd := NewNativeDetector(tt.display, tt.dummyMonitors, tt.minMonitors)

			if nd.display != tt.expectedDisplay {
				t.Errorf("Expected display %s, got %s", tt.expectedDisplay, nd.display)
			}

			if nd.minMonitors != tt.expectedMinMonitors {
				t.Errorf("Expected minMonitors %d, got %d", tt.expectedMinMonitors, nd.minMonitors)
			}

			if len(nd.dummyMonitors) != len(tt.dummyMonitors) {
				t.Errorf("Expected %d dummy monitors, got %d", len(tt.dummyMonitors), len(nd.dummyMonitors))
			}
		})
	}
}

// Note: DetectMonitors() is not tested here because it requires a running X server.
// In a real testing environment, you would use a mock X server or integration tests.
func TestNativeDetector_InterfaceCompliance(t *testing.T) {
	nd := NewNativeDetector(":0", []string{"dummy1", "dummy2"}, 3)

	// This test ensures that NativeDetector implements MonitorDetector interface
	var _ MonitorDetector = nd
}
