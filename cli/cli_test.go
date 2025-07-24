package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_Parse_DefaultValues(t *testing.T) {
	cli := NewCLI()

	args, err := cli.Parse([]string{"i3-config-generator"})
	if err != nil {
		t.Fatalf("Failed to parse default args: %v", err)
	}

	// Should use default values
	if args.ConfigPath != "" {
		t.Errorf("Expected empty ConfigPath, got %s", args.ConfigPath)
	}

	if args.LayoutName != DefaultLayoutName {
		t.Errorf("Expected layout %s, got %s", DefaultLayoutName, args.LayoutName)
	}

	// OutputPath should be expanded to absolute path
	expectedOutput := getDefaultOutputPath()
	if args.OutputPath != expectedOutput {
		t.Errorf("Expected OutputPath %s, got %s", expectedOutput, args.OutputPath)
	}
}

func TestCLI_Parse_CustomValues(t *testing.T) {
	cli := NewCLI()

	testArgs := []string{
		"i3-config-generator",
		"--config", "/tmp/test-config.yaml",
		"--output", "/tmp/test-output",
		"--layout", "one_mon",
	}

	args, err := cli.Parse(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse custom args: %v", err)
	}

	if args.ConfigPath != "/tmp/test-config.yaml" {
		t.Errorf("Expected ConfigPath /tmp/test-config.yaml, got %s", args.ConfigPath)
	}

	if args.OutputPath != "/tmp/test-output" {
		t.Errorf("Expected OutputPath /tmp/test-output, got %s", args.OutputPath)
	}

	if args.LayoutName != "one_mon" {
		t.Errorf("Expected layout one_mon, got %s", args.LayoutName)
	}
}

func TestCLI_Parse_ShortFlags(t *testing.T) {
	cli := NewCLI()

	testArgs := []string{
		"i3-config-generator",
		"-c", "/tmp/short-config.yaml",
		"-o", "/tmp/short-output",
		"-l", "no_mon",
	}

	args, err := cli.Parse(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse short flags: %v", err)
	}

	if args.ConfigPath != "/tmp/short-config.yaml" {
		t.Errorf("Expected ConfigPath /tmp/short-config.yaml, got %s", args.ConfigPath)
	}

	if args.OutputPath != "/tmp/short-output" {
		t.Errorf("Expected OutputPath /tmp/short-output, got %s", args.OutputPath)
	}

	if args.LayoutName != "no_mon" {
		t.Errorf("Expected layout no_mon, got %s", args.LayoutName)
	}
}

func TestCLI_Parse_InvalidLayout(t *testing.T) {
	cli := NewCLI()

	testArgs := []string{
		"i3-config-generator",
		"--layout", "invalid_layout",
	}

	_, err := cli.Parse(testArgs)
	if err == nil {
		t.Error("Expected error for invalid layout")
	}

	if !strings.Contains(err.Error(), "invalid layout name") {
		t.Errorf("Expected error about invalid layout, got: %v", err)
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected func() string
		wantErr  bool
	}{
		{
			name:  "empty path",
			input: "",
			expected: func() string {
				return ""
			},
			wantErr: false,
		},
		{
			name:  "absolute path",
			input: "/tmp/test",
			expected: func() string {
				return "/tmp/test"
			},
			wantErr: false,
		},
		{
			name:  "relative path",
			input: "test/path",
			expected: func() string {
				cwd, _ := os.Getwd()
				return filepath.Join(cwd, "test/path")
			},
			wantErr: false,
		},
		{
			name:  "home directory expansion",
			input: "~/test",
			expected: func() string {
				home, _ := os.UserHomeDir()
				return filepath.Join(home, "test")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandPath(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("expandPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expected := tt.expected()
				if result != expected {
					t.Errorf("expandPath() = %v, want %v", result, expected)
				}
			}
		})
	}
}

func TestIsValidLayout(t *testing.T) {
	tests := []struct {
		layout string
		valid  bool
	}{
		{"two_mon", true},
		{"one_mon", true},
		{"no_mon", true},
		{"invalid", false},
		{"", false},
		{"THREE_MON", false},
	}

	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			result := isValidLayout(tt.layout)
			if result != tt.valid {
				t.Errorf("isValidLayout(%s) = %v, want %v", tt.layout, result, tt.valid)
			}
		})
	}
}

func TestGetDefaultOutputPath(t *testing.T) {
	result := getDefaultOutputPath()

	// Should end with .i3/config
	if !strings.HasSuffix(result, ".i3/config") {
		t.Errorf("Expected default output path to end with .i3/config, got %s", result)
	}

	// Should be an absolute path
	if !filepath.IsAbs(result) {
		t.Errorf("Expected absolute path, got %s", result)
	}
}