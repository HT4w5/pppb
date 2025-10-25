package service

import (
	"fmt"
	"log"

	"github.com/HT4w5/pppb/internal/config"
	"github.com/HT4w5/pppb/internal/model"
)

type Service struct {
	cfg    *config.Config
	tasks  []*model.PPPTask
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

}

func (s *Service) runAll() {
	s.logger.Println("Starting all tasks.")

	results := make(chan model.PPPResult)

	for _, task := range s.tasks {
		go runTask(*task, results)
	}

	s.logger.Println("Waiting for tasks to finish...")
	successCount := 0

	for i := 0; i < len(s.tasks); i++ {
		result := <-results
		if result.Success {
			fmt.Printf("Task '%s' launched and completed successfully.\n", result.Tag)
			successCount++
		} else {
			fmt.Printf("Task '%s' failed. Error: %v\n", result.Tag, result.Error)
		}
	}

	s.logger.Println("All tasks ended. %d/%d successful.", len(s.tasks), successCount)
}
