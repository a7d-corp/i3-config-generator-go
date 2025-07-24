package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	Version           = "1.0.0"
	DefaultLayoutName = "two_mon"
)

// Args represents the parsed command-line arguments
type Args struct {
	ConfigPath  string
	OutputPath  string
	LayoutName  string
	ShowVersion bool
	ShowHelp    bool
}

// CLI handles command-line interface operations
type CLI struct {
	flagSet *flag.FlagSet
	args    *Args
}

// NewCLI creates a new CLI handler
func NewCLI() *CLI {
	args := &Args{}
	flagSet := flag.NewFlagSet("i3-config-generator", flag.ExitOnError)

	// Configuration file location flag
	flagSet.StringVar(&args.ConfigPath, "config", "",
		"Path to configuration file (default: $HOME/.config/i3-config-generator/config.yaml)")
	flagSet.StringVar(&args.ConfigPath, "c", "",
		"Path to configuration file (shorthand)")

	// Output file location flag
	defaultOutput := getDefaultOutputPath()
	flagSet.StringVar(&args.OutputPath, "output", defaultOutput,
		"Output file path for generated i3 configuration")
	flagSet.StringVar(&args.OutputPath, "o", defaultOutput,
		"Output file path for generated i3 configuration (shorthand)")

	// Layout selection flag
	flagSet.StringVar(&args.LayoutName, "layout", DefaultLayoutName,
		"Screen layout to use (two_mon, one_mon, no_mon)")
	flagSet.StringVar(&args.LayoutName, "l", DefaultLayoutName,
		"Screen layout to use (shorthand)")

	// Version flag
	flagSet.BoolVar(&args.ShowVersion, "version", false,
		"Show version information")
	flagSet.BoolVar(&args.ShowVersion, "v", false,
		"Show version information (shorthand)")

	// Help flag
	flagSet.BoolVar(&args.ShowHelp, "help", false,
		"Show help information")
	flagSet.BoolVar(&args.ShowHelp, "h", false,
		"Show help information (shorthand)")

	return &CLI{
		flagSet: flagSet,
		args:    args,
	}
}

// Parse parses the command-line arguments
func (cli *CLI) Parse(osArgs []string) (*Args, error) {
	// Custom usage function
	cli.flagSet.Usage = cli.printUsage

	if err := cli.flagSet.Parse(osArgs[1:]); err != nil {
		return nil, err
	}

	// Handle special flags first
	if cli.args.ShowVersion {
		cli.printVersion()
		os.Exit(0)
	}

	if cli.args.ShowHelp {
		cli.printUsage()
		os.Exit(0)
	}

	// Validate layout name
	if !isValidLayout(cli.args.LayoutName) {
		return nil, fmt.Errorf("invalid layout name: %s (valid options: two_mon, one_mon, no_mon)", cli.args.LayoutName)
	}

	// Expand paths
	if cli.args.ConfigPath != "" {
		expanded, err := expandPath(cli.args.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("invalid config path: %w", err)
		}
		cli.args.ConfigPath = expanded
	}

	expanded, err := expandPath(cli.args.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("invalid output path: %w", err)
	}
	cli.args.OutputPath = expanded

	return cli.args, nil
}

// printVersion prints version information
func (cli *CLI) printVersion() {
	fmt.Printf("i3-config-generator version %s\n", Version)
	fmt.Println("A dynamic i3 window manager configuration generator")
}

// printUsage prints usage information
func (cli *CLI) printUsage() {
	fmt.Printf("i3-config-generator v%s - Dynamic i3 configuration generator\n\n", Version)
	fmt.Println("USAGE:")
	fmt.Printf("  %s [OPTIONS]\n\n", os.Args[0])

	fmt.Println("DESCRIPTION:")
	fmt.Println("  Generates i3 window manager configuration files based on detected monitors")
	fmt.Println("  and user-defined layouts. Supports multiple screen configurations and")
	fmt.Println("  automatically maps workspaces to the appropriate displays.")
	fmt.Println()

	fmt.Println("OPTIONS:")
	cli.flagSet.PrintDefaults()
	fmt.Println()

	fmt.Println("EXAMPLES:")
	fmt.Printf("  # Generate config with default settings\n")
	fmt.Printf("  %s\n\n", os.Args[0])

	fmt.Printf("  # Generate config for single monitor setup\n")
	fmt.Printf("  %s --layout one_mon\n\n", os.Args[0])

	fmt.Printf("  # Use custom config and output locations\n")
	fmt.Printf("  %s --config ~/my-config.yaml --output ~/my-i3-config\n\n", os.Args[0])

	fmt.Printf("  # Generate config for no external monitors\n")
	fmt.Printf("  %s -l no_mon -o ~/.i3/laptop-config\n\n", os.Args[0])

	fmt.Println("LAYOUTS:")
	fmt.Println("  two_mon  - Two external monitors + laptop screen (default)")
	fmt.Println("  one_mon  - One external monitor + laptop screen")
	fmt.Println("  no_mon   - Laptop screen only")
	fmt.Println()

	fmt.Println("CONFIGURATION:")
	fmt.Println("  Default config location: $HOME/.config/i3-config-generator/config.yaml")
	fmt.Println("  Default output location: $HOME/.i3/config")
}

// getDefaultOutputPath returns the default output path for i3 configuration
func getDefaultOutputPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".i3/config" // Fallback
	}
	return filepath.Join(homeDir, ".i3", "config")
}

// expandPath expands ~ and environment variables in file paths
func expandPath(path string) (string, error) {
	if path == "" {
		return path, nil
	}

	// Expand ~ to home directory
	if path[:1] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Expand environment variables
	path = os.ExpandEnv(path)

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// isValidLayout checks if the layout name is valid
func isValidLayout(layout string) bool {
	validLayouts := map[string]bool{
		"two_mon": true,
		"one_mon": true,
		"no_mon":  true,
	}
	return validLayouts[layout]
}
