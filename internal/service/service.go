package service

import (
	"log"

	"github.com/HT4w5/pppcm/internal/config"
	"github.com/HT4w5/pppcm/internal/model"
)

type Service struct {
	cfg    *config.Config
	links  map[string]*model.PPPLink
	logger *log.Logger
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
	// Convert LinkConfigs to PPPTasks
	s.links = map[string]*model.PPPLink{}

	for _, c := range s.cfg.Links {
		s.links[c.Tag] = &model.PPPLink{
			Tag:    c.Tag,
			Task:   c.ToPPPTask(),
			Up:     false,
			IFName: c.IFName,
		}
	}
}

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
