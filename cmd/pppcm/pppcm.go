package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/HT4w5/pppcm/internal/config"
	"github.com/HT4w5/pppcm/internal/service"
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
	fmt.Println(VersionLine())
	cfg := config.New()
	if err := cfg.Load("config.json"); err != nil {
		log.Printf("[main] error loading config: %v\n", err)
		return
	}

	logger := log.Default()
	svc := service.New(cfg, logger)

	logger.Printf("[main] stopping existing links\n")
	svc.CheckAllLinks()
	svc.StopAllLinks()
	time.Sleep(5 * time.Second)

	if !cfg.Health.Enabled {
		// One shot start
		logger.Printf("[main] health disabled\n")
		svc.StartAllPPPTasks()
		logger.Printf("[main] exiting...\n")
		return
	}

	logger.Printf("[main] health enabled\n")
	for !svc.CheckAndRestart() {
		logger.Printf("[main] checking in 10s\n")
		time.Sleep(10 * time.Second)
	}
	logger.Printf("[main] start success. Checking in %ds\n", cfg.Health.CheckInterval)

	checkInterval := time.Duration(cfg.Health.CheckInterval) * time.Second
	checkTicker := time.NewTicker(checkInterval)
	done := make(chan os.Signal, 1)

	if !cfg.Health.ForceRestart {
		go func() {
			for {
				select {
				case <-done:
					return
				case <-checkTicker.C:
					for !svc.CheckAndRestart() {
						logger.Printf("[main] checking in 10s\n")
						time.Sleep(10 * time.Second)
					}
					logger.Printf("[main] check success. Waiting %ds\n", cfg.Health.CheckInterval)
					checkTicker.Reset(checkInterval)
				}
			}
		}()
	} else {
		forceRestartInterval := time.Duration(cfg.Health.ForceRestartInterval) * time.Second
		forceRestartTicker := time.NewTicker(forceRestartInterval)
		go func() {
			restart := time.Now()
			for {
				select {
				case <-done:
					return
				case <-checkTicker.C:
					for !svc.CheckAndRestart() {
						restart = time.Now()
						logger.Printf("[main] checking in 10s\n")
						time.Sleep(10 * time.Second)
					}
					logger.Printf("[main] check success. Waiting %ds\n", cfg.Health.CheckInterval)
					checkTicker.Reset(checkInterval)
				case <-forceRestartTicker.C:
					if time.Since(restart) > forceRestartInterval {
						logger.Printf("[main] force restart\n")
						svc.StartAllPPPTasks()
						for !svc.CheckAndRestart() {
							restart = time.Now()
							logger.Printf("[main] checking in 10s\n")
							time.Sleep(10 * time.Second)
						}
						checkTicker.Reset(checkInterval)
						forceRestartTicker.Reset(forceRestartInterval)
						logger.Printf("[main] force restart success. Waiting %ds\n", cfg.Health.ForceRestartInterval)
					}
				}
			}
		}()
	}

	signal.Notify(done, os.Interrupt)
	<-done
	logger.Printf("[main] exiting...\n")
}
