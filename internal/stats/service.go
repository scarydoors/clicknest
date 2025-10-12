package stats

import (
	"context"
	"log/slog"
	"time"
)

type Service struct {
	storage Storage
	logger *slog.Logger
}

type Storage interface {
	GetPageviews(ctx context.Context) (Timeseries, error)
}

func NewService(storage Storage, logger *slog.Logger) *Service {
	return &Service{
		storage: storage,
		logger: logger,
	}
}

type TimeseriesPoint struct {
	Timestamp time.Time
	Value uint64
}

type Timeseries []TimeseriesPoint

func (s *Service) GetPageviews(ctx context.Context) (Timeseries, error) {
	return s.storage.GetPageviews(ctx)
}
