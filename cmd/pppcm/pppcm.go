package main

import (
	"fmt"
	"runtime"

	"github.com/HT4w5/pppcm/internal/config"
)

const (
	VersionX = 0
	VersionY = 0
	VersionZ = 1
	Build    = "unknown"
)

func Version() string {
	return fmt.Sprintf("%v.%v.%v", VersionX, VersionY, VersionZ)
}

func VersionLine() string {
	return fmt.Sprintf("pppcm %s %s (%s %s/%s)", Version(), Build, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func main() {
	cfg := config.New()
	if err := cfg.Load("config.json"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(VersionLine())

	//svc := service.New(cfg)

	//sch := gocron.NewScheduler()
}
