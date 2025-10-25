package main

import (
	"github.com/HT4w5/pppb/internal/config"
	"github.com/HT4w5/pppb/internal/service"
)

func main() {
	cfg := config.New()
	cfg.Load("/etc/pppb/config.json")

	svc := service.New(cfg)
	svc.RunAll()
}
