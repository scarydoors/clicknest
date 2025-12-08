package stats

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	storage Storage
	logger *slog.Logger
	validate *validator.Validate 
}

type Storage interface {
	GetPageviews(ctx context.Context, params GetTimeseriesParameters) (Timeseries, error)
}

func NewService(storage Storage, logger *slog.Logger, validate *validator.Validate) *Service {
	return &Service{
		storage: storage,
		logger: logger,
		validate: validate,
	}
}

type TimeseriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value uint64 `json:"value"`
}

type Timeseries []TimeseriesPoint

type GetTimeseriesParameters struct {
	Domain string `validate:"required"`
	StartDate time.Time
	EndDate time.Time
	Interval time.Duration `validate:"interval_granularity=StartDate~EndDate:1000"`
}

func (s *Service) GetTimeseries(ctx context.Context, params GetTimeseriesParameters) (Timeseries, error) {
	if err := s.validate.Struct(params); err != nil {
		return Timeseries{}, err;
	}

	s.logger.Info("GetTimeseries", slog.Any("params", params))
	return s.storage.GetPageviews(ctx, params)
}
