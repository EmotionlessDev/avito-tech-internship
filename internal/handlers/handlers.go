package handlers

import (
	"log/slog"

	"github.com/EmotionlessDev/avito-tech-internship/internal/config"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
)

type Handlers struct {
	Health *HealthCheckHandler
	Team   *TeamHandler
}

func NewHandlers(
	logger *slog.Logger,
	cfg config.ConfigProvider,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	errorResponder *ErrorResponder,
) *Handlers {
	return &Handlers{
		Health: NewHealthCheckHandler(logger, cfg),
		Team:   NewTeamHandler(teamRepo, teamMemberRepo, errorResponder, logger),
	}
}
