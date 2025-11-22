package handlers

import (
	"log/slog"

	"github.com/EmotionlessDev/avito-tech-internship/internal/config"
	"github.com/EmotionlessDev/avito-tech-internship/internal/services"
)

type Handlers struct {
	Health *HealthCheckHandler
	Team   *TeamHandler
}

func NewHandlers(
	logger *slog.Logger,
	cfg config.ConfigProvider,
	errorResponder *ErrorResponder,
	teamService *services.TeamService,
) *Handlers {
	return &Handlers{
		Health: NewHealthCheckHandler(logger, cfg),
		Team:   NewTeamHandler(errorResponder, logger, teamService),
	}
}
