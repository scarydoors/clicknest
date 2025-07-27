package event

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/scarydoors/clicknest/internal/analytics"
)

type Service struct {
	storage Storage
	logger *slog.Logger
}

type Storage interface {
	InsertEvent(context.Context, analytics.Event) error
}

func NewService(storage Storage, logger *slog.Logger) *Service {
	return &Service{
		storage: storage,
		logger: logger,
	}	
}

func (s *Service) IngestEvent(ctx context.Context, event analytics.Event) error {
	if err := s.storage.InsertEvent(ctx, event); err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	return nil
}
