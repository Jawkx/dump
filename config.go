package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	FileStart string `toml:"file_start"`
	FileEnd   string `toml:"file_end"`
	CodeStart string `toml:"code_start"`
	CodeEnd   string `toml:"code_end"`
}

func NewConfig() *Config {
	return &Config{}
}

type ErrConfigFileNotFound struct {
	Path string
}

func (e ErrConfigFileNotFound) Error() string {
	return fmt.Sprintf("config file not found: %s", e.Path)
}

func (c *Config) Load(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return ErrConfigFileNotFound{Path: path}
	}

	_, err = toml.DecodeFile(path, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
		return err
	}
	return nil
}
