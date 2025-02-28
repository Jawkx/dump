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
	versionFlag := flag.Bool("version", false, "Display the version of the dump utility")
	versionShortFlag := flag.Bool("v", false, "Display the version of the dump utility")

	ignoreFlag := flag.String(
		"ignore",
		"",
		"A comma-separated list of patterns to ignore (e.g., \"*.log,temp/,config.ini\")",
	)
	ignoreShortFlag := flag.String(
		"i",
		"",
		"A comma-separated list of patterns to ignore (e.g., \"*.log,temp/,config.ini\")",
	)

	includeHiddenFlag := flag.Bool(
		"hidden",
		false,
		"Include hidden files and directories (those starting with a dot)",
	)

	helpFlag := flag.Bool("help", false, "Show help information")
	helpShortFlag := flag.Bool("h", false, "Show help information")

	flag.Parse()

	if *helpFlag || *helpShortFlag {
		printHelp()
		return
	}

	if *versionFlag || *versionShortFlag {
		fmt.Println("dump version:", version)
		return
	}

	paths := flag.Args()

	if len(paths) == 0 {
		fmt.Println("Error: No file or directory paths specified.")
		printHelp()
		os.Exit(1)
	}

	var ignorePatterns []string

	if *ignoreFlag != "" && *ignoreShortFlag != "" {
		fmt.Println("Error: Please use only one ignore flag, either -ignore or -i.")
		printHelp()
		os.Exit(1)
	} else if *ignoreFlag != "" {
		ignorePatterns = parseIgnorePatterns(*ignoreFlag)
	} else if *ignoreShortFlag != "" {
		ignorePatterns = parseIgnorePatterns(*ignoreShortFlag)
	}

	for _, path := range paths {
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
				processPath(match, ignorePatterns, *includeHiddenFlag)
			}
		} else {
			processPath(path, ignorePatterns, *includeHiddenFlag)
		}
	}
}

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

func shouldIgnore(path string, isDir bool, ignorePatterns []string, includeHidden bool) bool {
	if !includeHidden {
		baseName := filepath.Base(path)
		if len(baseName) > 0 && baseName[0] == '.' {
			return true
		}
	}

	if len(ignorePatterns) == 0 {
		return false
	}

	path = filepath.ToSlash(path)

	testPath := path
	if isDir && !strings.HasSuffix(testPath, "/") {
		testPath += "/"
	}

	for _, pattern := range ignorePatterns {
		pattern = filepath.ToSlash(pattern)

		if pattern == "." || pattern == "./" {
			if path == "." || path == "./" {
				return true
			}
		}

		isDirectoryPattern := strings.HasSuffix(pattern, "/")
		if isDirectoryPattern && isDir {
			dirPattern := strings.TrimSuffix(pattern, "/")

			if matched, _ := filepath.Match(dirPattern, strings.TrimSuffix(testPath, "/")); matched {
				return true
			}

			if strings.HasPrefix(path, pattern) {
				return true
			}
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched && !isDir {
			return true
		}

		matched, err = filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}

		if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") ||
			strings.Contains(pattern, "[") {
			if ok, _ := filepath.Match(pattern, path); ok {
				return true
			}

			if strings.Contains(pattern, "*/") {
				parts := strings.Split(pattern, "*/")
				if len(parts) >= 2 && strings.HasPrefix(path, parts[0]) {
					restPattern := strings.Join(parts[1:], "*/")
					restPath := strings.TrimPrefix(path, parts[0])

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

func processPath(path string, ignorePatterns []string, includeHidden bool) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error accessing '%s': %v\n", path, err)
		return
	}

	if shouldIgnore(path, info.IsDir(), ignorePatterns, includeHidden) {
		return
	}

	if info.IsDir() {
		processDirectory(path, ignorePatterns, includeHidden)
	} else {
		dumpFile(path)
	}
}

func processDirectory(dirPath string, ignorePatterns []string, includeHidden bool) {
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing '%s': %v\n", path, err)
			return nil
		}

		isDir := d.IsDir()

		// Check if the item should be ignored
		if shouldIgnore(path, isDir, ignorePatterns, includeHidden) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		if !isDir {
			dumpFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory '%s': %v\n", dirPath, err)
	}
}

func dumpFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading '%s': %v\n", filePath, err)
		return
	}

	ext := filepath.Ext(filePath)
	if ext != "" {
		ext = ext[1:]
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

func containsGlobPattern(path string) bool {
	return strings.ContainsAny(path, "*?[")
}

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
	fmt.Println(
		"  -i=\"<pattern1>,<pattern2>,...\"  Comma-separated list of patterns to ignore",
	)
	fmt.Println("  -hidden        Include hidden files and directories (starting with a dot)")
	fmt.Println("  -help          Show this help information")
	fmt.Println("\nExamples:")
	fmt.Println("  dump file.txt                        # Dump a single file")
	fmt.Println(
		"  dump src/                            # Dump all files in a directory recursively",
	)
	fmt.Println("  dump *.go                            # Dump all Go files in current directory")
	fmt.Println("  dump -ignore=\"*.log,temp/\" project/   # Ignore .log files and temp directory")
	fmt.Println("  dump -hidden project/                # Include hidden files like .gitignore")
}
