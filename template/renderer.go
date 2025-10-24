package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/a7d-corp/i3-config-generator-go/config"
	"github.com/a7d-corp/i3-config-generator-go/monitor"
)

//go:embed i3.tmpl
var embeddedTemplates embed.FS

// TemplateData represents the data structure passed to the template
type TemplateData struct {
	I3                  config.I3Config
	Colors              config.ColorConfig
	Layout              ResolvedLayoutConfig
	ApplicationBindings map[string]string
	StartupPrograms     []string
	WindowOverrides     []string
	DetectedMonitors    *monitor.DetectedMonitors
}

// ResolvedLayoutConfig is a layout config with monitor references resolved
type ResolvedLayoutConfig struct {
	GapsInner          int
	GapsOuter          int
	MoveWorkspace      map[string]string // Resolved to actual monitor names
	WorkspaceToDisplay map[string]string // Resolved to actual monitor names
}

// Renderer handles template rendering operations
type Renderer struct {
	templateDir string
}

// NewRenderer creates a new template renderer
func NewRenderer(templateDir string) *Renderer {
	if templateDir == "" {
		// Default to template directory relative to current working directory
		templateDir = "template"
	}

	return &Renderer{
		templateDir: templateDir,
	}
}

// Render generates the i3 configuration by rendering the template with the given data
func (r *Renderer) Render(cfg *config.Config, layoutName string, detectedMonitors *monitor.DetectedMonitors) (string, error) {
	// Get the specified layout
	layout, err := cfg.GetLayout(layoutName)
	if err != nil {
		return "", fmt.Errorf("failed to get layout: %w", err)
	}

	// Resolve monitor references in the layout
	resolvedLayout, err := r.resolveLayoutReferences(layout, detectedMonitors)
	if err != nil {
		return "", fmt.Errorf("failed to resolve layout references: %w", err)
	}

	// Prepare template data
	templateData := &TemplateData{
		I3:                  cfg.I3,
		Colors:              cfg.Colors,
		Layout:              *resolvedLayout,
		ApplicationBindings: cfg.ApplicationBindings,
		StartupPrograms:     cfg.StartupPrograms,
		WindowOverrides:     cfg.WindowOverrides,
		DetectedMonitors:    detectedMonitors,
	}

	// Load and render template
	return r.renderTemplate("i3.tmpl", templateData)
}

// resolveLayoutReferences converts layout role references to actual monitor names
func (r *Renderer) resolveLayoutReferences(layout *config.LayoutConfig, detectedMonitors *monitor.DetectedMonitors) (*ResolvedLayoutConfig, error) {
	resolved := &ResolvedLayoutConfig{
		GapsInner:          layout.GapsInner,
		GapsOuter:          layout.GapsOuter,
		MoveWorkspace:      make(map[string]string),
		WorkspaceToDisplay: make(map[string]string),
	}

	// Resolve MoveWorkspace references
	for keybind, role := range layout.MoveWorkspace {
		monitorName := detectedMonitors.GetMonitorByRole(role)
		if monitorName == "" {
			return nil, fmt.Errorf("unknown monitor role: %s", role)
		}
		resolved.MoveWorkspace[keybind] = monitorName
	}

	// Resolve WorkspaceToDisplay references
	for workspace, role := range layout.WorkspaceToDisplay {
		monitorName := detectedMonitors.GetMonitorByRole(role)
		if monitorName == "" {
			return nil, fmt.Errorf("unknown monitor role: %s", role)
		}
		resolved.WorkspaceToDisplay[workspace] = monitorName
	}

	return resolved, nil
}

// renderTemplate loads and renders the specified template file
func (r *Renderer) renderTemplate(templateFile string, data *TemplateData) (string, error) {
	var tmpl *template.Template
	var err error

	// First, try to load from filesystem (for development/customization)
	if r.templateDir != "" {
		templatePath := filepath.Join(r.templateDir, templateFile)
		if _, statErr := os.Stat(templatePath); statErr == nil {
			// External template file exists, use it
			tmpl, err = template.ParseFiles(templatePath)
			if err != nil {
				return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
			}
		}
	}

	// If no external template found, use embedded template
	if tmpl == nil {
		templateContent, err := embeddedTemplates.ReadFile(templateFile)
		if err != nil {
			return "", fmt.Errorf("template file not found: %s (neither in filesystem nor embedded)", templateFile)
		}
		tmpl, err = template.New(templateFile).Parse(string(templateContent))
		if err != nil {
			return "", fmt.Errorf("failed to parse embedded template %s: %w", templateFile, err)
		}
	}

	// Create a buffer to capture the rendered output
	var output []byte
	buf := &writeBuffer{data: &output}

	// Execute the template
	if err := tmpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return string(output), nil
}

// writeBuffer is a simple buffer that implements io.Writer
type writeBuffer struct {
	data *[]byte
}

func (w *writeBuffer) Write(p []byte) (n int, err error) {
	*w.data = append(*w.data, p...)
	return len(p), nil
}

// RenderToFile renders the template and writes the output to a file
func (r *Renderer) RenderToFile(cfg *config.Config, layoutName string, detectedMonitors *monitor.DetectedMonitors, outputPath string) error {
	content, err := r.Render(cfg, layoutName, detectedMonitors)
	if err != nil {
		return err
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write the rendered configuration to file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
