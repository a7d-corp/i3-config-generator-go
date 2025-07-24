package main

import (
	"fmt"
	"log"
	"os"

	"github.com/a7d-corp/i3-config-generator-go/cli"
	"github.com/a7d-corp/i3-config-generator-go/config"
	"github.com/a7d-corp/i3-config-generator-go/monitor"
	"github.com/a7d-corp/i3-config-generator-go/template"
)

func main() {
	// Parse command-line arguments
	cliHandler := cli.NewCLI()
	args, err := cliHandler.Parse(os.Args)
	if err != nil {
		log.Fatalf("Error parsing arguments: %v", err)
	}

	// Load configuration
	loader := config.NewLoader("")
	var cfg *config.Config
	if args.ConfigPath != "" {
		cfg, err = loader.LoadFromFile(args.ConfigPath)
	} else {
		cfg, err = loader.Load()
	}
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("✓ Configuration loaded successfully\n")

	// Detect monitors if enabled
	var detectedMonitors *monitor.DetectedMonitors
	if cfg.UseDetectedMonitors {
		fmt.Printf("✓ Detecting monitors...\n")
		detector, err := cfg.CreateDetector()
		if err != nil {
			log.Fatalf("Failed to create monitor detector: %v", err)
		}
		detectedMonitors, err = detector.DetectMonitors()
		if err != nil {
			log.Fatalf("Failed to detect monitors: %v", err)
		}
		fmt.Printf("✓ Detected %d monitors: %s\n", len(detectedMonitors.All), detectedMonitors.Primary)
		if len(detectedMonitors.All) > 1 {
			fmt.Printf("  - Primary: %s, Left: %s, Right: %s\n",
				detectedMonitors.Primary, detectedMonitors.Left, detectedMonitors.Right)
		}
	} else {
		fmt.Printf("✓ Using static monitor configuration\n")
	}

	// Render the template
	fmt.Printf("✓ Rendering i3 configuration for layout: %s\n", args.LayoutName)
	renderer := template.NewRenderer("")

	renderedConfig, err := renderer.Render(cfg, args.LayoutName, detectedMonitors)
	if err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}

	// Write the configuration to the output file
	fmt.Printf("✓ Writing configuration to: %s\n", args.OutputPath)
	if err := renderer.RenderToFile(cfg, args.LayoutName, detectedMonitors, args.OutputPath); err != nil {
		log.Fatalf("Failed to write configuration file: %v", err)
	}

	// Show summary
	fmt.Printf("\n🎉 i3 configuration generated successfully!\n")
	fmt.Printf("   Layout: %s\n", args.LayoutName)
	fmt.Printf("   Output: %s\n", args.OutputPath)
	if detectedMonitors != nil {
		fmt.Printf("   Monitors: %d detected\n", len(detectedMonitors.All))
	}
	fmt.Printf("   Size: %.1f KB\n", float64(len(renderedConfig))/1024)

	fmt.Println("\nTo use this configuration:")
	fmt.Printf("   1. Backup your current i3 config (if any)\n")
	fmt.Printf("   2. Restart i3: i3-msg restart\n")
	fmt.Printf("   3. Or reload config: i3-msg reload\n")
}
