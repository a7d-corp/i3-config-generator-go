package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a7d-corp/i3-config-generator-go/config"
	"github.com/a7d-corp/i3-config-generator-go/monitor"
)

func TestRenderer_resolveLayoutReferences(t *testing.T) {
	renderer := NewRenderer("")

	detectedMonitors := &monitor.DetectedMonitors{
		Primary: "eDP-1",
		Left:    "HDMI-1",
		Right:   "DP-1",
		All:     []string{"eDP-1", "HDMI-1", "DP-1"},
	}

	layout := &config.LayoutConfig{
		GapsInner: 20,
		GapsOuter: 0,
		MoveWorkspace: map[string]string{
			"Ctrl+Shift+1": "left_display",
			"Ctrl+Shift+2": "right_display",
			"Ctrl+Shift+3": "primary_display",
		},
		WorkspaceToDisplay: map[string]string{
			"1": "left_display",
			"2": "left_display",
			"3": "right_display",
			"4": "primary_display",
		},
	}

	resolved, err := renderer.resolveLayoutReferences(layout, detectedMonitors)
	if err != nil {
		t.Fatalf("Failed to resolve layout references: %v", err)
	}

	// Check basic properties
	if resolved.GapsInner != 20 {
		t.Errorf("Expected GapsInner 20, got %d", resolved.GapsInner)
	}

	// Check MoveWorkspace resolution
	expectedMoveWorkspace := map[string]string{
		"Ctrl+Shift+1": "HDMI-1",
		"Ctrl+Shift+2": "DP-1",
		"Ctrl+Shift+3": "eDP-1",
	}

	for keybind, expectedMonitor := range expectedMoveWorkspace {
		if actualMonitor := resolved.MoveWorkspace[keybind]; actualMonitor != expectedMonitor {
			t.Errorf("MoveWorkspace[%s]: expected %s, got %s", keybind, expectedMonitor, actualMonitor)
		}
	}

	// Check WorkspaceToDisplay resolution
	expectedWorkspaceToDisplay := map[string]string{
		"1": "HDMI-1",
		"2": "HDMI-1",
		"3": "DP-1",
		"4": "eDP-1",
	}

	for workspace, expectedMonitor := range expectedWorkspaceToDisplay {
		if actualMonitor := resolved.WorkspaceToDisplay[workspace]; actualMonitor != expectedMonitor {
			t.Errorf("WorkspaceToDisplay[%s]: expected %s, got %s", workspace, expectedMonitor, actualMonitor)
		}
	}
}

func TestRenderer_resolveLayoutReferences_InvalidRole(t *testing.T) {
	renderer := NewRenderer("")

	detectedMonitors := &monitor.DetectedMonitors{
		Primary: "eDP-1",
		Left:    "HDMI-1",
		Right:   "DP-1",
	}

	layout := &config.LayoutConfig{
		MoveWorkspace: map[string]string{
			"Ctrl+Shift+1": "invalid_role",
		},
	}

	_, err := renderer.resolveLayoutReferences(layout, detectedMonitors)
	if err == nil {
		t.Error("Expected error for invalid monitor role")
	}

	if !strings.Contains(err.Error(), "unknown monitor role") {
		t.Errorf("Expected error message about unknown monitor role, got: %v", err)
	}
}

func TestRenderer_renderTemplate(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple test template
	testTemplate := `# Test i3 config
set $mod {{.I3.ModKey}}

# Colors
set $base00 {{.Colors.Base00}}

# Gaps
gaps inner {{.Layout.GapsInner}}
gaps outer {{.Layout.GapsOuter}}

# Application bindings
{{range $key, $value := .ApplicationBindings}}
assign {{$key}} {{$value}}
{{end}}

# Startup programs
{{range .StartupPrograms}}
exec --no-startup-id {{.}}
{{end}}

# Monitor setup
exec --no-startup-id xrandr --output {{.DetectedMonitors.Primary}} --primary
`

	templatePath := filepath.Join(tempDir, "i3.tmpl")
	if err := os.WriteFile(templatePath, []byte(testTemplate), 0644); err != nil {
		t.Fatalf("Failed to write test template: %v", err)
	}

	renderer := NewRenderer(tempDir)

	// Create test data
	templateData := &TemplateData{
		I3: config.I3Config{
			ModKey: "Mod4",
		},
		Colors: config.ColorConfig{
			Base00: "#1B2B34",
		},
		Layout: ResolvedLayoutConfig{
			GapsInner: 20,
			GapsOuter: 0,
		},
		ApplicationBindings: map[string]string{
			"[class=\"^Firefox$\"]":  "1",
			"[class=\"^Chromium$\"]": "3",
		},
		StartupPrograms: []string{
			"/usr/bin/numlockx on",
			"/usr/bin/compton -b",
		},
		DetectedMonitors: &monitor.DetectedMonitors{
			Primary: "eDP-1",
		},
	}

	result, err := renderer.renderTemplate("i3.tmpl", templateData)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Check that key elements are present in the output
	expectedElements := []string{
		"set $mod Mod4",
		"set $base00 #1B2B34",
		"gaps inner 20",
		"gaps outer 0",
		"assign [class=\"^Firefox$\"] 1",
		"assign [class=\"^Chromium$\"] 3",
		"exec --no-startup-id /usr/bin/numlockx on",
		"exec --no-startup-id /usr/bin/compton -b",
		"exec --no-startup-id xrandr --output eDP-1 --primary",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected output to contain '%s', but it was missing from:\n%s", expected, result)
		}
	}
}

func TestRenderer_renderTemplate_FileNotFound(t *testing.T) {
	renderer := NewRenderer("/nonexistent/path")

	templateData := &TemplateData{}

	_, err := renderer.renderTemplate("i3.tmpl", templateData)
	if err == nil {
		t.Error("Expected error when template file doesn't exist")
	}

	if !strings.Contains(err.Error(), "template file not found") {
		t.Errorf("Expected error about template file not found, got: %v", err)
	}
}
