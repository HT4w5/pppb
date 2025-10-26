package service

import (
	"log"

	"github.com/HT4w5/pppb/internal/config"
	"github.com/HT4w5/pppb/internal/model"
)

type Service struct {
	cfg         *config.Config
	connections map[string]*model.PPPConnection
	logger      *log.Logger
}

func New(cfg *config.Config) *Service {
	s := &Service{
		cfg:    cfg,
		logger: log.Default(),
	}
	s.init()
	return s
}

func (s *Service) init() {
	// Convert ConnectionConfigs to PPPTasks
	s.connections = map[string]*model.PPPConnection{}

	for _, c := range s.cfg.Connections {
		s.connections[c.Tag] = &model.PPPConnection{
			Tag:     c.Tag,
			Task:    c.ToPPPTask(),
			Up:      false,
			Healthy: false,
			IFName:  c.IFName,
		}
	}
}

func (s *Service) RunAllPPPTasks() {
	s.logger.Println("[main] Starting all ppp tasks")

	results := make(chan model.PPPResult)
	startSignal := make(chan struct{})

	// Spawn tasks
	for _, c := range s.connections {
		go runPPPTask(*c.Task, startSignal, results)
	}
	// Send start signal
	close(startSignal)

	s.logger.Println("[main] Waiting for ppp tasks to finish...")
	successCount := 0

	for i := 0; i < len(s.connections); i++ {
		result := <-results
		s.connections[result.Tag].Up = result.Success
		if result.Success {
			s.logger.Printf("[%s] started successfully\n", result.Tag)
			successCount++
		} else {
			s.logger.Printf("[%s] failed to start: %v\n", result.Tag, result.Error)
		}
	}
	s.logger.Printf("[main] All ppp tasks ended. %d/%d successful", successCount, len(s.connections))
}
