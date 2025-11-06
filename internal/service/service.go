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

func New(cfg *config.Config, l *log.Logger) *Service {
	s := &Service{
		cfg:    cfg,
		logger: l,
	}
	s.init()
	return s
}

func (s *Service) init() {
	// Convert LinkConfigs to PPPTasks
	s.links = make(map[string]*model.PPPLink, len(s.cfg.Links))

	for _, c := range s.cfg.Links {
		s.links[c.Tag] = &model.PPPLink{
			Tag:    c.Tag,
			Task:   c.ToPPPTask(),
			Up:     false,
			IFName: c.IFName,
		}
	}
}

// Check links, if less than expected, restart all, return false.
// If equal or more, return true.
func (s *Service) CheckAndRestart() bool {
	upCount := s.CheckAllLinks()
	if upCount >= s.cfg.Health.Expected {
		s.logger.Printf("[service] %d links up, satisfies expected %d\n", upCount, s.cfg.Health.Expected)
		return true
	}
	s.logger.Printf("[service] %d links up, less than expected %d\n", upCount, s.cfg.Health.Expected)
	s.StopAllLinks()
	s.StartAllPPPTasks()
	return false
}
