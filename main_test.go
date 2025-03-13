package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseIgnorePatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Empty string", "", nil},
		{"Single pattern", "*.log", []string{"*.log"}},
		{"Multiple patterns", "*.log,temp/,config.ini", []string{"*.log", "temp/", "config.ini"}},
		{"With spaces", " *.log , temp/ , config.ini ", []string{"*.log", "temp/", "config.ini"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseIgnorePatterns(tc.input)

			if tc.expected == nil && result != nil {
				t.Fatalf("Expected nil, got %v", result)
			}

			if tc.expected != nil {
				if len(tc.expected) != len(result) {
					t.Fatalf("Expected %v, got %v", tc.expected, result)
				}

				for i, v := range tc.expected {
					if v != result[i] {
						t.Errorf("Expected %s at position %d, got %s", v, i, result[i])
					}
				}
			}
		})
	}
}

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		isDir         bool
		patterns      []string
		includeHidden bool
		shouldIgn     bool
	}{
		{"No patterns", "file.txt", false, nil, false, false},
		{"Simple file match", "file.log", false, []string{"*.log"}, false, true},
		{"File doesn't match", "file.txt", false, []string{"*.log"}, false, false},
		{"Directory match", "temp", true, []string{"temp/"}, false, true},
		{"Subdirectory match", "logs/temp/file.txt", true, []string{"logs/temp/"}, false, true},
		{"Current directory", ".", true, []string{"."}, false, true},
		{"Path with wildcard", "logs/2023/error.log", false, []string{"logs/*/*.log"}, false, true},
		{"Exact file match", "config.ini", false, []string{"config.ini"}, false, true},
		{"Complex pattern", "src/temp/cache.tmp", false, []string{"**/temp/**"}, false, true},
		{"Hidden file ignored", ".gitignore", false, nil, false, true},
		{"Hidden file included", ".gitignore", false, nil, true, false},
		{"Hidden dir ignored", ".git", true, nil, false, true},
		{"Hidden dir included", ".git", true, nil, true, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := shouldIgnore(tc.path, tc.isDir, tc.patterns, tc.includeHidden)
			if result != tc.shouldIgn {
				t.Errorf(
					"Expected shouldIgnore to return %v for path '%s' with patterns %v and includeHidden=%v, got %v",
					tc.shouldIgn,
					tc.path,
					tc.patterns,
					tc.includeHidden,
					result,
				)
			}
		})
	}
}

func TestContainsGlobPattern(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"file.txt", false},
		{"*.txt", true},
		{"file?.txt", true},
		{"file[1-3].txt", true},
		{"dir/subdir/file.txt", false},
		{"dir/*/file.txt", true},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			result := containsGlobPattern(tc.path)
			if result != tc.expected {
				t.Errorf("Expected containsGlobPattern('%s') to return %v, got %v",
					tc.path, tc.expected, result)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.FileStart != "--- FILE-START: {{.FilePath}} ---" {
		t.Errorf("Expected default FileStart, got %s", config.FileStart)
	}

	if config.FileEnd != "--- FILE-END ---" {
		t.Errorf("Expected default FileEnd, got %s", config.FileEnd)
	}

	if config.ContentStart != "``` {{.FileExt}}" {
		t.Errorf("Expected default ContentStart, got %s", config.ContentStart)
	}

	if config.ContentEnd != "```" {
		t.Errorf("Expected default ContentEnd, got %s", config.ContentEnd)
	}
}

func TestConfigLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.toml")
	configContent := `
file_start = "== FILE: {{.FilePath}} =="
file_end = "== END FILE =="
code_start = "<< {{.FileExt}} >>"
code_end = "<< END >>"
`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Test loading config
	config := NewConfig()
	err = config.Load(configPath)

	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.FileStart != "== FILE: {{.FilePath}} ==" {
		t.Errorf("Expected custom FileStart, got %s", config.FileStart)
	}

	if config.FileEnd != "== END FILE ==" {
		t.Errorf("Expected custom FileEnd, got %s", config.FileEnd)
	}

	if config.ContentStart != "<< {{.FileExt}} >>" {
		t.Errorf("Expected custom ContentStart, got %s", config.ContentStart)
	}

	if config.ContentEnd != "<< END >>" {
		t.Errorf("Expected custom ContentEnd, got %s", config.ContentEnd)
	}

	// Test loading non-existent config file
	config = NewConfig()
	err = config.Load(filepath.Join(tmpDir, "nonexistent.toml"))

	if err != nil {
		t.Errorf("Loading non-existent file should return nil, got: %v", err)
	}
}

func TestConfigLoadFromPaths(t *testing.T) {
	// Create a temporary config file
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath1 := filepath.Join(tmpDir, "config1.toml")
	configPath2 := filepath.Join(tmpDir, "config2.toml")

	configContent := `
file_start = "== FILE: {{.FilePath}} =="
file_end = "== END FILE =="
code_start = "<< {{.FileExt}} >>"
code_end = "<< END >>"
`

	err = os.WriteFile(configPath1, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	err = os.WriteFile(configPath2, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Test loading from multiple paths (first doesn't exist, second does)
	config := NewConfig()
	err = config.LoadFromPaths([]string{configPath1, configPath2})

	if err != nil {
		t.Fatalf("LoadFromPaths failed: %v", err)
	}

	if config.FileStart != "== FILE: {{.FilePath}} ==" {
		t.Errorf("Expected custom FileStart from second path, got %s", config.FileStart)
	}
}

func TestDumpFile(t *testing.T) {
	// Create a temporary file
	content := "test content"
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create a custom config
	config := &Config{
		FileStart:    "START: {{.FilePath}}",
		FileEnd:      "END",
		ContentStart: "```",
		ContentEnd:   "```",
	}

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function we're testing
	dumpFile(tmpfile.Name(), config)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf = make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Verify output contains our content
	if output == "" {
		t.Error("Expected output, got nothing")
	}

	if !strings.Contains(output, content) {
		t.Errorf("Output doesn't contain file content: %s", output)
	}

	// Verify custom format is used
	if !strings.Contains(output, "START: "+tmpfile.Name()) {
		t.Errorf("Output doesn't use custom format: %s", output)
	}
}

// TestProcessPath simulates processing both files and directories
func TestProcessPath(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "dump-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file in the temp directory
	testFilePath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	config := NewConfig()

	// Test processing a file
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	processPath(testFilePath, nil, false, config)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf = make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Verify file was processed
	if output == "" {
		t.Error("Expected output for file processing, got nothing")
	}

	// Now test with a directory
	r, w, _ = os.Pipe()
	os.Stdout = w

	processPath(tempDir, nil, false, config)

	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	buf = make([]byte, 1024)
	n, _ = r.Read(buf)
	output = string(buf[:n])

	// Verify directory was processed
	if output == "" {
		t.Error("Expected output for directory processing, got nothing")
	}
}
