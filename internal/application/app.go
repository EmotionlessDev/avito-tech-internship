package application

import (
	"database/sql"
	"log/slog"
)

type Config interface {
}

type Application struct {
	Config Config
	Logger *slog.Logger
	DB     *sql.DB
}

func New(cfg Config, logger *slog.Logger, db *sql.DB) *Application {
	return &Application{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}
}
