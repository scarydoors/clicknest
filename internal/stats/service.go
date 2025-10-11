package stats

import (
	"log/slog"
	"time"
)

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

type PageviewEntry struct {
	time time.Time
	value uint64
}

type PageviewResult []PageviewEntry

func (s *Service) GetPageviews() {
}
