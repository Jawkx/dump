package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	FileStart    string `toml:"file_start"`
	FileEnd      string `toml:"file_end"`
	ContentStart string `toml:"code_start"`
	ContentEnd   string `toml:"code_end"`
}

func NewConfig() *Config {
	return &Config{
		FileStart:    "--- FILE-START: {{.FilePath}} ---",
		FileEnd:      "--- FILE-END ---",
		ContentStart: "``` {{.FileExt}}",
		ContentEnd:   "```",
	}
}

func (c *Config) Load(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	_, err = toml.DecodeFile(path, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
		return err
	}
	return nil
}

func (c *Config) LoadFromPaths(paths []string) error {
	for _, path := range paths {
		err := c.Load(path)
		if err == nil {
			return nil
		}
		if !os.IsNotExist(err) {
			return err
		}
		// TODO: Implement better logging system
		// fmt.Printf("Config file not found at %s, trying next path...\n", path)
	}
	return fmt.Errorf("no config file found in provided paths")
}
