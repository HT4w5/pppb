package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	pppdCommand = "pppd"
)

type Config struct {
	Sessions []*SessionConfig `yaml:"sessions"`
}

type SessionConfig struct {
	Name   string   `yaml:"name"`
	Comand string   `yaml:"command"`
	Args   []string `yaml:"args"`
}

func New() *Config {
	return &Config{
		[]*SessionConfig{},
	}
}

func (c *Config) Load(config string) {
	configBytes, err := os.ReadFile(config)
	if err != nil {
		log.Panicf("Failed to open config: %s, %v", config, err)
	}

	if err := yaml.Unmarshal(configBytes, &c); err != nil {
		log.Panicf("Failed to parse config: %s, %v", config, err)
	}

	// Use default pppd command
	for _, s := range c.Sessions {
		s.Comand = pppdCommand
	}
}
