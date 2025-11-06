package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HT4w5/pppcm/internal/model"
)

type Config struct {
	Links  []*LinkConfig `json:"links"`
	RunDir string        `json:"run_dir"`
}

func New() *Config {
	return &Config{
		Links:  []*LinkConfig{},
		RunDir: "/var/run",
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

	if err := c.validate(); err != nil {
		return fmt.Errorf("Failed to validate config: %s, %w", config, err)
	}

	return nil
}

func (c *Config) validate() error {
	// Runtime variable directory
	_, err := os.ReadDir(c.RunDir)
	if err != nil {
		return fmt.Errorf("Can't open runtime variable directory %s: %w", c.RunDir, err)
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
