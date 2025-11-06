package service

import (
	"os"
	"os/exec"

	"github.com/HT4w5/pppcm/internal/model"
)

func (s *Service) StartAllPPPTasks() {
	s.logger.Println("[main] Starting all pppd tasks")

	results := make(chan model.PPPResult)
	startSignal := make(chan struct{})

	// Spawn tasks
	for _, c := range s.links {
		go runPPPTask(*c.Task, startSignal, results)
	}
	// Send start signal
	close(startSignal)

	s.logger.Println("[main] Waiting for pppd tasks to finish...")
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
	s.logger.Printf("[main] All pppd tasks ended. %d/%d successful", successCount, len(s.links))
}

func (s *Service) CheckAllLinks() {
	s.logger.Println("[main] Checking all ppp links")
}

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
