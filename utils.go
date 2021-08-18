package main

import (
	"flag"
	"fmt"
	"os"
)

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, bool, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	backend := flag.Bool("backend", false, "run in backend mode")

	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", *backend, err
	}

	return configPath, *backend, nil
}
