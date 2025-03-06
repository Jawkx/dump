package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

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
	fmt.Printf("``` %s\n", filePath)
	fmt.Print(string(content))
	if !strings.HasSuffix(string(content), "\n") {
		fmt.Println()
	}
	fmt.Println("```")
	fmt.Println("---FILE-END---")
	fmt.Println()
}
