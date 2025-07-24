package monitor

import (
	"fmt"
	"sort"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

// NativeDetector implements monitor detection using native X11 RandR bindings
type NativeDetector struct {
	display       string
	dummyMonitors []string
	minMonitors   int
}

// NewNativeDetector creates a new native detector
func NewNativeDetector(display string, dummyMonitors []string, minMonitors int) *NativeDetector {
	if display == "" {
		display = ":0" // Default X display
	}
	if minMonitors == 0 {
		minMonitors = 3 // Default minimum
	}

	return &NativeDetector{
		display:       display,
		dummyMonitors: dummyMonitors,
		minMonitors:   minMonitors,
	}
}

// DetectMonitors detects connected monitors using native X11 RandR calls
func (nd *NativeDetector) DetectMonitors() (*DetectedMonitors, error) {
	// Connect to X server
	conn, err := xgb.NewConnDisplay(nd.display)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to X display %s: %w", nd.display, err)
	}
	defer conn.Close()

	// Initialize RandR extension
	err = randr.Init(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RandR extension: %w", err)
	}

	// Get root window
	setup := xproto.Setup(conn)
	root := setup.DefaultScreen(conn).Root

	// Get screen resources
	resources, err := randr.GetScreenResources(conn, root).Reply()
	if err != nil {
		return nil, fmt.Errorf("failed to get screen resources: %w", err)
	}

	var connectedOutputs []string
	var primaryOutput string

	// Query each output to check if it's connected
	for _, output := range resources.Outputs {
		outputInfo, err := randr.GetOutputInfo(conn, output, resources.ConfigTimestamp).Reply()
		if err != nil {
			continue // Skip on error
		}

		// Check if output is connected
		if outputInfo.Connection == randr.ConnectionConnected {
			name := string(outputInfo.Name)
			connectedOutputs = append(connectedOutputs, name)

			// Check if this is the primary output
			if outputInfo.Crtc != 0 {
				crtcInfo, err := randr.GetCrtcInfo(conn, outputInfo.Crtc, resources.ConfigTimestamp).Reply()
				if err == nil && len(crtcInfo.Outputs) > 0 {
					// Check if this CRTC is the primary
					primary, err := randr.GetOutputPrimary(conn, root).Reply()
					if err == nil && primary.Output == output {
						primaryOutput = name
					}
				}
			}
		}
	}

	// If no primary output was detected, use the first connected output
	if primaryOutput == "" && len(connectedOutputs) > 0 {
		primaryOutput = connectedOutputs[0]
	}

	// Sort outputs for consistent ordering
	sort.Strings(connectedOutputs)

	// Pad with dummy monitors if needed
	paddedOutputs := nd.padWithDummyMonitors(connectedOutputs)

	// Assign roles - same logic as before
	var left, right string
	if len(paddedOutputs) >= 2 {
		left = paddedOutputs[1] // Second monitor as left
	}
	if len(paddedOutputs) >= 3 {
		right = paddedOutputs[2] // Third monitor as right
	}

	// Ensure primary is set to first monitor if not detected
	if primaryOutput == "" && len(paddedOutputs) > 0 {
		primaryOutput = paddedOutputs[0]
	}

	return &DetectedMonitors{
		Primary: primaryOutput,
		Left:    left,
		Right:   right,
		All:     paddedOutputs,
	}, nil
}

// padWithDummyMonitors pads the monitor list with dummy monitors to meet minimum requirements
func (nd *NativeDetector) padWithDummyMonitors(monitors []string) []string {
	result := make([]string, len(monitors))
	copy(result, monitors)

	// Add dummy monitors until we reach the minimum
	dummyIndex := 0
	for len(result) < nd.minMonitors && dummyIndex < len(nd.dummyMonitors) {
		result = append(result, nd.dummyMonitors[dummyIndex])
		dummyIndex++
	}

	return result
}

// GetDisplayInfo returns information about the X display connection
func (nd *NativeDetector) GetDisplayInfo() (string, error) {
	conn, err := xgb.NewConnDisplay(nd.display)
	if err != nil {
		return "", fmt.Errorf("failed to connect to X display %s: %w", nd.display, err)
	}
	defer conn.Close()

	setup := xproto.Setup(conn)
	screen := setup.DefaultScreen(conn)

	return fmt.Sprintf("Display: %s, Screen: %dx%d",
		nd.display,
		screen.WidthInPixels,
		screen.HeightInPixels), nil
}
