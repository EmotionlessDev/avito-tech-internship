package application

import (
	"database/sql"
	"log/slog"

	"github.com/EmotionlessDev/avito-tech-internship/internal/config"
	"github.com/EmotionlessDev/avito-tech-internship/internal/handlers"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
)

type Application struct {
	Config   config.ConfigProvider
	Logger   *slog.Logger
	DB       *sql.DB
	Handlers *handlers.Handlers
}

func New(
	cfg config.ConfigProvider,
	logger *slog.Logger,
	db *sql.DB,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	errorResponder *handlers.ErrorResponder,
) *Application {
	return &Application{
		Config:   cfg,
		Logger:   logger,
		DB:       db,
		Handlers: handlers.NewHandlers(logger, cfg, teamRepo, teamMemberRepo, errorResponder),
	}
}
