package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HT4w5/pppb/internal/model"
)

type Config struct {
	Connections []*ConnectionConfig `json:"connections"`
}

func New() *Config {
	return &Config{
		[]*ConnectionConfig{},
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
	for _, cc := range c.Connections {
		if len(cc.Tag) == 0 {
			return fmt.Errorf("Connection with no tag present")
		}
		if len(cc.TTYName) == 0 || len(cc.User) == 0 || len(cc.Password) == 0 || len(cc.IFName) == 0 {
			return fmt.Errorf("Invalid connection: %s", cc.Tag)
		}
	}
	return nil
}

type ConnectionConfig struct {
	Tag      string `json:"tag"`
	TTYName  string `json:"ttyname"`
	User     string `json:"user"`
	Password string `json:"password"`
	IFName   string `json:"ifname"`
}

func (c *ConnectionConfig) toPPPTask() *model.PPPTask {
	task := &model.PPPTask{
		Comand: model.PPPDefaultCommand,
		Tag:    c.Tag,
		Args:   c.getArgs(),
	}
	return task
}

func (c *ConnectionConfig) getArgs() []string {
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
