package application

import (
	"database/sql"
	"log/slog"

	"github.com/EmotionlessDev/avito-tech-internship/internal/config"
	"github.com/EmotionlessDev/avito-tech-internship/internal/handlers"
)

type Application struct {
	Config   config.ConfigProvider
	Logger   *slog.Logger
	DB       *sql.DB
	Handlers *handlers.Handlers
}

func New(cfg config.ConfigProvider, logger *slog.Logger, db *sql.DB) *Application {
	return &Application{
		Config:   cfg,
		Logger:   logger,
		DB:       db,
		Handlers: handlers.NewHandlers(logger, cfg),
	}
}
