package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var version string

func main() {

	if version == "" {
		version = "development"
	}

	config := NewConfig()

	configPaths := []string{
		filepath.Join(userHomeDir(), ".config/dump.toml"),
		filepath.Join(userHomeDir(), ".config/dump/config.toml"),
	}

	err := config.LoadFromPaths(configPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

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
				processPath(match, ignorePatterns, *includeHiddenFlag, config)
			}
		} else {
			processPath(path, ignorePatterns, *includeHiddenFlag, config)
		}
	}
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

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Unable to determine user home directory: %v\n", err)
		return ""
	}
	return home
}
