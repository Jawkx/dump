// dump - A Simple File and Directory Content Dumper
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const version = "1.0.0"

func main() {
	// Define command-line flags
	versionFlag := flag.Bool("version", false, "Display the version of the dump utility")
	ignoreFlag := flag.String(
		"ignore",
		"",
		"A comma-separated list of patterns to ignore (e.g., \"*.log,temp/,config.ini\")",
	)
	helpFlag := flag.Bool("help", false, "Show help information")

	// Parse command-line flags
	flag.Parse()

	// Show help if requested
	if *helpFlag {
		printHelp()
		return
	}

	// Display version and exit if version flag is set
	if *versionFlag {
		fmt.Println("dump version:", version)
		return
	}

	// Collect paths to process
	paths := flag.Args()

	// Check if no paths are provided
	if len(paths) == 0 {
		fmt.Println("Error: No file or directory paths specified.")
		printHelp()
		os.Exit(1)
	}

	// Parse ignore patterns
	ignorePatterns := parseIgnorePatterns(*ignoreFlag)

	// Process each path
	for _, path := range paths {
		// Check if the path is a glob pattern
		if containsGlobPattern(path) {
			matches, err := filepath.Glob(path)
			if err != nil {
				fmt.Printf("Error with glob pattern '%s': %v\n", path, err)
				continue
			}

			if len(matches) == 0 {
				fmt.Printf("Warning: No files matched the pattern '%s'\n", path)
				continue
			}

			for _, match := range matches {
				processPath(match, ignorePatterns)
			}
		} else {
			processPath(path, ignorePatterns)
		}
	}
}

// Parse comma-separated ignore patterns
func parseIgnorePatterns(ignoreStr string) []string {
	if ignoreStr == "" {
		return nil
	}

	patterns := strings.Split(ignoreStr, ",")
	for i := range patterns {
		patterns[i] = strings.TrimSpace(patterns[i])
	}
	return patterns
}

// Check if a path should be ignored based on patterns
func shouldIgnore(path string, isDir bool, ignorePatterns []string) bool {
	if len(ignorePatterns) == 0 {
		return false
	}

	// Convert path to use forward slashes for consistency
	path = filepath.ToSlash(path)

	// If it's a directory, add a trailing slash for pattern matching
	testPath := path
	if isDir && !strings.HasSuffix(testPath, "/") {
		testPath += "/"
	}

	for _, pattern := range ignorePatterns {
		// Normalize patterns to use forward slashes
		pattern = filepath.ToSlash(pattern)

		// Check if pattern represents the current directory
		if pattern == "." || pattern == "./" {
			if path == "." || path == "./" {
				return true
			}
		}

		// Directory pattern handling
		isDirectoryPattern := strings.HasSuffix(pattern, "/")
		if isDirectoryPattern && isDir {
			// Remove trailing slash for matching
			dirPattern := strings.TrimSuffix(pattern, "/")

			// Match exact directory or parent directory
			if matched, _ := filepath.Match(dirPattern, strings.TrimSuffix(testPath, "/")); matched {
				return true
			}

			// Check if it's a subdirectory of an ignored directory
			if strings.HasPrefix(path, pattern) {
				return true
			}
		}

		// File pattern handling
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched && !isDir {
			return true
		}

		// Path pattern handling
		matched, err = filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}

		// Handle glob patterns that might match the full path
		if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") ||
			strings.Contains(pattern, "[") {
			if ok, _ := filepath.Match(pattern, path); ok {
				return true
			}

			// Special handling for patterns like "logs/*/*.log"
			if strings.Contains(pattern, "*/") {
				parts := strings.Split(pattern, "*/")
				if len(parts) >= 2 && strings.HasPrefix(path, parts[0]) {
					restPattern := strings.Join(parts[1:], "*/")
					restPath := strings.TrimPrefix(path, parts[0])

					// Check if any subdirectory matches the pattern
					dirs := strings.Split(restPath, "/")
					for i := range dirs {
						subPath := strings.Join(dirs[i:], "/")
						if ok, _ := filepath.Match(restPattern, subPath); ok {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

// Process a file or directory path
func processPath(path string, ignorePatterns []string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error accessing '%s': %v\n", path, err)
		return
	}

	// Check if the item should be ignored
	if shouldIgnore(path, info.IsDir(), ignorePatterns) {
		return
	}

	if info.IsDir() {
		// Process directory recursively
		processDirectory(path, ignorePatterns)
	} else {
		// Process individual file
		dumpFile(path)
	}
}

// Process a directory recursively
func processDirectory(dirPath string, ignorePatterns []string) {
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing '%s': %v\n", path, err)
			return nil
		}

		isDir := d.IsDir()

		// Check if the item should be ignored
		if shouldIgnore(path, isDir, ignorePatterns) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		// Dump file content if it's a regular file
		if !isDir {
			dumpFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory '%s': %v\n", dirPath, err)
	}
}

// Dump the content of a file
func dumpFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading '%s': %v\n", filePath, err)
		return
	}

	// Determine file extension for the code block format
	ext := filepath.Ext(filePath)
	if ext != "" {
		ext = ext[1:] // Remove the dot
	}

	fmt.Println("---FILE-START---")
	fmt.Printf("``` %s\n", filepath.Base(filePath))
	fmt.Print(string(content))
	if !strings.HasSuffix(string(content), "\n") {
		fmt.Println()
	}
	fmt.Println("```")
	fmt.Println("---FILE-END---")
	fmt.Println()
}

// Check if a path contains glob patterns
func containsGlobPattern(path string) bool {
	return strings.ContainsAny(path, "*?[")
}

// Print help information
func printHelp() {
	fmt.Println("dump - A Simple File and Directory Content Dumper")
	fmt.Println("\nUsage:")
	fmt.Println("  dump [options] <file_path1> <file_path2> ...")
	fmt.Println("\nArguments:")
	fmt.Println("  <file_path>    Path to a file or directory. Supports glob patterns.")
	fmt.Println("\nOptions:")
	fmt.Println("  -version       Display the version of the dump utility")
	fmt.Println(
		"  -ignore=\"<pattern1>,<pattern2>,...\"  Comma-separated list of patterns to ignore",
	)
	fmt.Println("  -help          Show this help information")
	fmt.Println("\nExamples:")
	fmt.Println("  dump file.txt                        # Dump a single file")
	fmt.Println(
		"  dump src/                            # Dump all files in a directory recursively",
	)
	fmt.Println("  dump *.go                            # Dump all Go files in current directory")
	fmt.Println("  dump -ignore=\"*.log,temp/\" project/   # Ignore .log files and temp directory")
}
