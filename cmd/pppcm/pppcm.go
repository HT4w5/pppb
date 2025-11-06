package main

import (
	"fmt"

	"github.com/HT4w5/pppcm/internal/config"
	"github.com/HT4w5/pppcm/internal/service"
)

func main() {
	cfg := config.New()
	if err := cfg.Load("config.json"); err != nil {
		fmt.Println(err)
		return
	}

	svc := service.New(cfg)
	//svc.StartAllPPPTasks()
	svc.CheckAllLinks()
}
