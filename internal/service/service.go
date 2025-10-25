package service

import "github.com/HT4w5/pppb/internal/config"

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	s := &Service{
		cfg: cfg,
	}
	s.init()
	return s
}

func (s *Service) init() {

}
