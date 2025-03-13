package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func processPath(path string, ignorePatterns []string, includeHidden bool, config *Config) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error accessing '%s': %v\n", path, err)
		return
	}

	if shouldIgnore(path, info.IsDir(), ignorePatterns, includeHidden) {
		return
	}

	if info.IsDir() {
		processDirectory(path, ignorePatterns, includeHidden, config)
	} else {
		dumpFile(path, config)
	}
}

func processDirectory(dirPath string, ignorePatterns []string, includeHidden bool, config *Config) {
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing '%s': %v\n", path, err)
			return nil
		}

		isDir := d.IsDir()

		if shouldIgnore(path, isDir, ignorePatterns, includeHidden) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		if !isDir {
			dumpFile(path, config)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory '%s': %v\n", dirPath, err)
	}
}

type FileData struct {
	FilePath string
	Content  string
	FileExt  string
}

func dumpFile(filePath string, config *Config) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading '%s': %v\n", filePath, err)
		return
	}

	ext := filepath.Ext(filePath)
	if ext != "" {
		// Remove the leading dot
		ext = ext[1:]
	}

	data := FileData{
		FilePath: filePath,
		Content:  string(content),
		FileExt:  ext,
	}

	fileStartString, err := parseTemplate(config.FileStart, data)
	if err != nil {
		fmt.Printf("Error parsing file start string")
		return
	}

	contentStartString, err := parseTemplate(config.ContentStart, data)
	if err != nil {
		fmt.Printf("Error parsing file start string")
		return
	}

	contentEndString, err := parseTemplate(config.ContentEnd, data)
	if err != nil {
		fmt.Printf("Error parsing file start string")
		return
	}

	fileEndString, err := parseTemplate(config.FileEnd, data)
	if err != nil {
		fmt.Printf("Error parsing file start string")
		return
	}

	var sb strings.Builder
	sb.WriteString(fileStartString + "\n")
	sb.WriteString(contentStartString + "\n")
	sb.WriteString(string(content) + "\n")
	sb.WriteString(contentEndString + "\n")
	sb.WriteString(fileEndString + "\n")

	fmt.Println(sb.String())
}

func parseTemplate(input string, data FileData) (string, error) {
	tmpl, err := template.New("").Parse(input)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
