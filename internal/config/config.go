package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Connections []*ConnectionConfig `json:"connections"`
}

type ConnectionConfig struct {
	Tag      string `json:"tag"`
	TTYName  string `json:"ttyname"`
	User     string `json:"user"`
	Password string `json:"password"`
	IFName   string `json:"ifname"`
}

func New() *Config {
	return &Config{
		[]*ConnectionConfig{},
	}
}

func (c *Config) Load(config string) {
	configBytes, err := os.ReadFile(config)
	if err != nil {
		log.Panicf("Failed to open config: %s, %v", config, err)
	}

	if err := json.Unmarshal(configBytes, &c); err != nil {
		log.Panicf("Failed to parse config: %s, %v", config, err)
	}
}
