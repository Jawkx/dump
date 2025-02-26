package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	redColor   = "\033[31m"
	resetColor = "\033[0m"
)

func dumpContext(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fmt.Printf("``` %s\n%s\n```\n", filePath, string(content))
	return nil
}

func processFile(filePath string) error {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			redColor+"Error getting absolute path: %s: %s\n"+resetColor,
			filePath,
			err,
		)
		return err
	}

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		fmt.Fprintf(
			os.Stderr,
			redColor+"Error: File does not exist: %s\n"+resetColor,
			absFilePath,
		)
		return err
	}

	if err := dumpContext(absFilePath); err != nil {
		fmt.Fprintf(
			os.Stderr,
			redColor+"Error dumping context for %s: %s\n"+resetColor,
			absFilePath,
			err,
		)
		return err
	}
	return nil
}

func processPath(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			redColor+"Error accessing path %s: %s\n"+resetColor,
			path,
			err,
		)
		return err
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					redColor+"Error accessing %s: %s\n"+resetColor,
					filePath,
					err,
				)
				return err
			}
			if !fileInfo.IsDir() {
				return processFile(filePath)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		return processFile(path)
	}
	return nil
}

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dump [options] <file_path1> <file_path2> ...\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if versionFlag {
		fmt.Println("dump version 0.0.1")
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintf(
			os.Stderr,
			redColor+"Error: Please provide at least one file path.\n"+resetColor,
		)
		flag.Usage()
		os.Exit(1)
	}

	for _, path := range flag.Args() {
		if err := processPath(path); err != nil {
			continue
		}
	}
}

func expandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, path[1:]), nil
}
