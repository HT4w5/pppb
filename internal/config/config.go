package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/HT4w5/pppcm/internal/model"
)

type Config struct {
	Links  []*LinkConfig `json:"links"`  // Link configs
	Health *HealthConfig `json:"health"` // Link health check config
}

func New() *Config {
	return &Config{
		Links: []*LinkConfig{},
		Health: &HealthConfig{
			RunDir:               "/var/run",
			Enabled:              false, // Don't check link status by default
			CheckInterval:        300,   // 5m
			ForceRestart:         false,
			ForceRestartInterval: 86400, // 1d
		},
	}
}

func (c *Config) Load(config string) error {
	configBytes, err := os.ReadFile(config)
	if err != nil {
		return fmt.Errorf("Failed to open config: %s, %w", config, err)
	}

	if err := json.Unmarshal(configBytes, &c); err != nil {
		return fmt.Errorf("Failed to parse config: %s, %w", config, err)
	}

	// Strip suffix
	c.Health.RunDir = strings.TrimSuffix(c.Health.RunDir, "/")

	if err := c.validate(); err != nil {
		return fmt.Errorf("Failed to validate config: %s, %w", config, err)
	}

	return nil
}

func (c *Config) validate() error {
	if c.Health.Enabled {
		// Runtime variable directory
		_, err := os.ReadDir(c.Health.RunDir)
		if err != nil {
			return fmt.Errorf("Can't open runtime variable directory %s: %w", c.Health.RunDir, err)
		}

		// Expected
		if c.Health.Expected > len(c.Links) {
			return fmt.Errorf("Expected (%d) greater than number of links (%d)", c.Health.Expected, len(c.Links))
		}

		// Intervals
		if c.Health.CheckInterval < 0 {
			return fmt.Errorf("Invalid CheckInterval %v, must be equal to or greater than 300", c.Health.CheckInterval)
		}

		if c.Health.ForceRestart {
			if c.Health.ForceRestartInterval < c.Health.CheckInterval {
				return fmt.Errorf("Invalid ForceRestartInterval %v, must be equal to or greater than CheckInterval (%v)", c.Health.ForceRestartInterval, c.Health.CheckInterval)
			}
		}
	}

	// Links
	for _, cc := range c.Links {
		if len(cc.Tag) == 0 {
			return fmt.Errorf("Link with no tag present")
		}
		if len(cc.TTYName) == 0 || len(cc.User) == 0 || len(cc.Password) == 0 || len(cc.IFName) == 0 {
			return fmt.Errorf("Invalid link: %s", cc.Tag)
		}
	}
	return nil
}

type LinkConfig struct {
	Tag      string `json:"tag"`
	TTYName  string `json:"ttyname"`
	User     string `json:"user"`
	Password string `json:"password"`
	IFName   string `json:"ifname"`
}

func (c *LinkConfig) ToPPPTask() *model.PPPTask {
	task := &model.PPPTask{
		Comand: model.PPPDefaultCommand,
		Tag:    c.Tag,
		Args:   c.getArgs(),
	}
	return task
}

func (c *LinkConfig) getArgs() []string {
	return []string{
		"plugin",
		"pppoe.so",
		c.TTYName,
		"user",
		c.User,
		"password",
		c.Password,
		"ifname",
		c.IFName,
	}
}

type HealthConfig struct {
	RunDir               string `json:"run_dir"`                // Runtime variable directory where [ifname].pid exists
	Enabled              bool   `json:"enabled"`                // Enable health check
	Expected             int    `json:"expected"`               // Expected number of links (restart if less)
	CheckInterval        int    `json:"check_interval"`         // Check interval in seconds
	ForceRestart         bool   `json:"force_restart"`          // Enable force restart
	ForceRestartInterval int    `json:"force_restart_interval"` // Force restart interval in seconds
}
