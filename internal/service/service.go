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
