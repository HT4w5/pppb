package service

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/HT4w5/pppcm/internal/model"
)

// Concurrently run pppd to start all links.
// Returns number of successful calls to pppd.
func (s *Service) StartAllPPPTasks() int {
	s.logger.Println("[service] starting all pppd tasks")

	results := make(chan model.PPPResult)
	startSignal := make(chan struct{})

	// Spawn tasks
	for _, c := range s.links {
		go runPPPTask(*c.Task, startSignal, results)
	}
	// Send start signal
	close(startSignal)

	s.logger.Println("[service] waiting for pppd tasks to finish...")
	successCount := 0

	for i := 0; i < len(s.links); i++ {
		result := <-results
		s.links[result.Tag].Up = result.Success
		if result.Success {
			s.logger.Printf("[%s] started successfully\n", result.Tag)
			successCount++
		} else {
			s.logger.Printf("[%s] failed to start: %v\n", result.Tag, result.Error)
		}
	}
	s.logger.Printf("[service] all pppd tasks ended. %d/%d successful\n", successCount, len(s.links))
	return successCount
}

// Look at pppd PID files to see whether links are up.
// Return number of links up.
func (s *Service) CheckAllLinks() int {
	s.logger.Println("[service] checking all ppp links")

	upCount := 0
	for k, v := range s.links {
		pidBytes, err := os.ReadFile(s.cfg.Daemon.RunDir + "/" + v.IFName + ".pid")
		if err == nil {
			pidStr := strings.TrimSpace(string(pidBytes))
			pid, err := strconv.Atoi(pidStr)
			if err == nil {
				s.logger.Printf("[%s] link up, pid %d, iface %s\n", k, pid, v.IFName)
				v.Up = true
				v.PID = pid
				upCount++
				continue
			}
		}
		s.logger.Printf("[%s] link down\n", k)
		v.Up = false
		v.PID = 0
	}
	s.logger.Printf("[service] %d/%d links up\n", upCount, len(s.links))
	return upCount
}

func (s *Service) StopAllLinks() int {
	s.logger.Println("[service] stopping all ppp links")

	stopCount := 0
	for k, v := range s.links {
		if !v.Up {
			s.logger.Printf("[%s] link already down, skipping\n", k)
			continue
		}

		p, err := os.FindProcess(v.PID)
		if err != nil {
			s.logger.Printf("[%s] invalid PID %d, skipping\n", k, v.PID)
			continue
		}

		// Validate pppd process
		cmdlineBytes, err := os.ReadFile("/proc/" + strconv.Itoa(v.PID) + "/cmdline")
		if !(cmdlineBytes[0] == 'p' && cmdlineBytes[1] == 'p' && cmdlineBytes[2] == 'p' && cmdlineBytes[3] == 'd') {
			s.logger.Printf("[%s] process cmdline mismatch at PID %d, skipping\n", k, v.PID)
			continue
		}

		if err := p.Signal(syscall.SIGTERM); err != nil {
			// If SIGTERM fails, try the forceful SIGKILL
			s.logger.Printf("[%s] SIGTERM failed for PID %d\n", k, v.PID)

			// SIGKILL is usually required to be os.Kill() for cross-platform compatibility
			if kerr := p.Kill(); kerr != nil {
				s.logger.Printf("[%s] failed to send SIGKILL to PID %d: %v\n", k, v.PID, kerr)
				continue
			}
			fmt.Printf("[%s] SIGKILL successfully sent to PID %d\n", k, v.PID)
		}
		s.logger.Printf("[%s] stopped link with PID %d\n", k, v.PID)
		stopCount++
	}
	s.logger.Printf("[service] stopped %d links\n", stopCount)
	return stopCount
}

// Runs pppd to start a link.
// Returns a PPPResult.
func runPPPTask(task model.PPPTask, startSignal <-chan struct{}, results chan<- model.PPPResult) {
	// Block until start signal
	<-startSignal

	cmd := exec.Command(task.Comand, task.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		results <- model.PPPResult{
			Tag:     task.Tag,
			Success: false,
			Error:   err,
		}
		return
	}

	results <- model.PPPResult{
		Tag:     task.Tag,
		Success: true,
		Error:   nil,
	}
}
