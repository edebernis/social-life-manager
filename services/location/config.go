package main

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Config struct for app configuration
type Config struct {
	SQL struct {
		ConnMaxIdleTime string `yaml:"connMaxIdleTime"`
		ConnMaxLifeTime string `yaml:"connMaxLifeTime"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
	} `yaml:"sql"`
	JWT struct {
		Algorithm string
	} `yaml:"jwt"`
}

func newConfig(r io.Reader) (*Config, error) {
	var config Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, fmt.Errorf("Failed to decode YAML config. %w", err)
	}

	return &config, nil
}
