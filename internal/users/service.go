package users

import "log/slog"

type Service struct {
	logger *slog.Logger
	storage Storage 
}

type Storage interface {
}

func NewService(storage Storage, logger *slog.Logger) *Service {
	return &Service{
		storage: storage,
		logger: logger,
	}
}
