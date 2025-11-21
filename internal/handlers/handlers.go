package handlers

import (
	"log/slog"

	"github.com/EmotionlessDev/avito-tech-internship/internal/config"
)

type Handlers struct {
	Health *HealthCheckHandler
}

func NewHandlers(logger *slog.Logger, cfg config.ConfigProvider) *Handlers {
	return &Handlers{
		Health: NewHealthCheckHandler(logger, cfg),
	}
}
