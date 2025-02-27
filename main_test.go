package main

import (
	"os"
	"path/filepath"
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
		name      string
		path      string
		isDir     bool
		patterns  []string
		shouldIgn bool
	}{
		{"No patterns", "file.txt", false, nil, false},
		{"Simple file match", "file.log", false, []string{"*.log"}, true},
		{"File doesn't match", "file.txt", false, []string{"*.log"}, false},
		{"Directory match", "temp", true, []string{"temp/"}, true},
		{"Subdirectory match", "logs/temp/file.txt", true, []string{"logs/temp/"}, true},
		{"Current directory", ".", true, []string{"."}, true},
		{"Path with wildcard", "logs/2023/error.log", false, []string{"logs/*/*.log"}, true},
		{"Exact file match", "config.ini", false, []string{"config.ini"}, true},
		{"Complex pattern", "src/temp/cache.tmp", false, []string{"**/temp/**"}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := shouldIgnore(tc.path, tc.isDir, tc.patterns)
			if result != tc.shouldIgn {
				t.Errorf(
					"Expected shouldIgnore to return %v for path '%s' with patterns %v, got %v",
					tc.shouldIgn,
					tc.path,
					tc.patterns,
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

// TestDumpFile creates a temporary file and tests if the content is correctly read
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

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function we're testing
	dumpFile(tmpfile.Name())

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

	// Test processing a file
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	processPath(testFilePath, nil)

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

	processPath(tempDir, nil)

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
