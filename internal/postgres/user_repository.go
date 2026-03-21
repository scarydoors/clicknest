package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	conn *pgx.Conn
	logger *slog.Logger
}

func NewUserRepository(conn *pgx.Conn, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		conn: conn,
		logger: logger,
	}
}
