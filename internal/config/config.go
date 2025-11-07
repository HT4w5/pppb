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
	Daemon *DaemonConfig `json:"daemon"` // Daemon config
}

func New() *Config {
	return &Config{
		Links: []*LinkConfig{},
		Daemon: &DaemonConfig{
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
	c.Daemon.RunDir = strings.TrimSuffix(c.Daemon.RunDir, "/")

	if err := c.validate(); err != nil {
		return fmt.Errorf("Failed to validate config: %s, %w", config, err)
	}

	return nil
}

func (c *Config) validate() error {
	if c.Daemon.Enabled {
		// Runtime variable directory
		_, err := os.ReadDir(c.Daemon.RunDir)
		if err != nil {
			return fmt.Errorf("can't open runtime variable directory %s: %w", c.Daemon.RunDir, err)
		}

		// Expected
		if c.Daemon.Expected > len(c.Links) {
			return fmt.Errorf("expected (%d) greater than number of links (%d)", c.Daemon.Expected, len(c.Links))
		}

		// Intervals
		if c.Daemon.CheckInterval < 300 {
			return fmt.Errorf("invalid CheckInterval %v, must be equal to or greater than 300", c.Daemon.CheckInterval)
		}

		if c.Daemon.ForceRestart {
			if c.Daemon.ForceRestartInterval < c.Daemon.CheckInterval {
				return fmt.Errorf("invalid ForceRestartInterval %v, must be equal to or greater than CheckInterval (%v)", c.Daemon.ForceRestartInterval, c.Daemon.CheckInterval)
			}
		}
	}

	// Links
	numLinks := len(c.Links)
	tags := make(map[string]struct{}, numLinks)
	ttynames := make(map[string]struct{}, numLinks)
	ifnames := make(map[string]struct{}, numLinks)
	for _, cc := range c.Links {
		if len(cc.Tag) == 0 {
			return fmt.Errorf("link with no tag present")
		}
		if len(cc.TTYName) == 0 || len(cc.User) == 0 || len(cc.Password) == 0 || len(cc.IFName) == 0 {
			return fmt.Errorf("invalid link: %s", cc.Tag)
		}
		tags[cc.Tag] = struct{}{}
		ttynames[cc.TTYName] = struct{}{}
		ifnames[cc.IFName] = struct{}{}
	}
	if len(tags) != numLinks {
		return fmt.Errorf("duplicate link tags")
	}
	if len(ttynames) != numLinks {
		return fmt.Errorf("duplicate link ttynames")
	}
	if len(ifnames) != numLinks {
		return fmt.Errorf("duplicate link ifnames")
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

type DaemonConfig struct {
	RunDir               string `json:"run_dir"`                // Runtime variable directory where [ifname].pid exists
	Enabled              bool   `json:"enabled"`                // Enable daemon
	Expected             int    `json:"expected"`               // Expected number of links (restart if less)
	CheckInterval        int    `json:"check_interval"`         // Check interval in seconds
	ForceRestart         bool   `json:"force_restart"`          // Enable force restart
	ForceRestartInterval int    `json:"force_restart_interval"` // Force restart interval in seconds
}
