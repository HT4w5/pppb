package main

import (
	"fmt"

	"github.com/HT4w5/pppb/internal/config"
	"github.com/HT4w5/pppb/internal/service"
)

func main() {
	cfg := config.New()
	if err := cfg.Load("config.json"); err != nil {
		fmt.Println(err)
		return
	}

	svc := service.New(cfg)
	svc.RunAll()
}
