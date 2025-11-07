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

	MinRetryInterval = 10   // 10s
	MaxRetryInterval = 1200 // 2m
	Banner           = `██████╗ ██████╗ ██████╗  ██████╗███╗   ███╗
██╔══██╗██╔══██╗██╔══██╗██╔════╝████╗ ████║
██████╔╝██████╔╝██████╔╝██║     ██╔████╔██║
██╔═══╝ ██╔═══╝ ██╔═══╝ ██║     ██║╚██╔╝██║
██║     ██║     ██║     ╚██████╗██║ ╚═╝ ██║
╚═╝     ╚═╝     ╚═╝      ╚═════╝╚═╝     ╚═╝`
)

func Version() string {
	return fmt.Sprintf("%v.%v.%v", VersionX, VersionY, VersionZ)
}

func VersionLine() string {
	return fmt.Sprintf("pppcm %s %s (%s %s/%s)", Version(), Build, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func main() {
	fmt.Println(Banner)
	fmt.Println(VersionLine())
	cfg := config.New()
	if err := cfg.Load("config.json"); err != nil {
		log.Printf("[main] error loading config: %v\n", err)
		return
	}

	logger := log.Default()
	svc := service.New(cfg, logger)

	// Stop existing links
	logger.Printf("[main] stopping existing links\n")
	svc.CheckAllLinks()
	svc.StopAllLinks()
	time.Sleep(5 * time.Second)

	// Daemon disabled
	if !cfg.Daemon.Enabled {
		// One shot start
		logger.Printf("[main] daemon disabled\n")
		svc.StartAllPPPTasks()
		logger.Printf("[main] exiting...\n")
		return
	}

	// Daemon enabled
	logger.Printf("[main] daemon enabled\n")
	// Start all links
	i := MinRetryInterval
	for !svc.CheckAndRestart() {
		logger.Printf("[main] (re)checking in %ds\n", i)
		time.Sleep(time.Duration(i) * time.Second)
		i = intervalFallback(i)
	}
	logger.Printf("[main] start success. Checking in %ds\n", cfg.Daemon.CheckInterval)

	// Set check ticker
	checkInterval := time.Duration(cfg.Daemon.CheckInterval) * time.Second
	checkTicker := time.NewTicker(checkInterval)
	done := make(chan os.Signal, 1)

	// No force restart
	if !cfg.Daemon.ForceRestart {
		go func() {
			for {
				select {
				case <-done:
					return
				case <-checkTicker.C:
					i := MinRetryInterval
					for !svc.CheckAndRestart() {
						logger.Printf("[main] (re)checking in %ds\n", i)
						time.Sleep(time.Duration(i) * time.Second)
						i = intervalFallback(i)
					}
					logger.Printf("[main] check success. Waiting %ds\n", cfg.Daemon.CheckInterval)
					checkTicker.Reset(checkInterval)
				}
			}
		}()
	} else { // With force restart
		forceRestartInterval := time.Duration(cfg.Daemon.ForceRestartInterval) * time.Second
		forceRestartTicker := time.NewTicker(forceRestartInterval)
		go func() {
			// Time of last restart
			restart := time.Now()
			for {
				select {
				case <-done:
					return
				case <-checkTicker.C:
					i := MinRetryInterval
					for !svc.CheckAndRestart() {
						restart = time.Now()
						logger.Printf("[main] (re)checking in %ds\n", i)
						time.Sleep(time.Duration(i) * time.Second)
						i = intervalFallback(i)
					}
					logger.Printf("[main] check success. Waiting %ds\n", cfg.Daemon.CheckInterval)
					checkTicker.Reset(checkInterval)
				case <-forceRestartTicker.C:
					if time.Since(restart) > forceRestartInterval {
						logger.Printf("[main] force restart\n")
						svc.StartAllPPPTasks()
						i := MinRetryInterval
						for !svc.CheckAndRestart() {
							restart = time.Now()
							logger.Printf("[main] (re)checking in %ds\n", i)
							time.Sleep(time.Duration(i) * time.Second)
							i = intervalFallback(i)
						}
						checkTicker.Reset(checkInterval)
						forceRestartTicker.Reset(forceRestartInterval)
						logger.Printf("[main] force restart success. Waiting %ds\n", cfg.Daemon.ForceRestartInterval)
					}
				}
			}
		}()
	}

	signal.Notify(done, os.Interrupt)
	<-done
	logger.Printf("[main] exiting...\n")
}

// Exponentially fallback retrial interval to prevent abuse.
func intervalFallback(i int) int {
	if i >= MaxRetryInterval {
		return i
	} else {
		return i * 2
	}
}
