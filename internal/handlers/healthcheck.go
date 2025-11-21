package handlers

import (
	"log/slog"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/config"
	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
)

type HealthCheckHandler struct {
	Logger *slog.Logger
	Config config.ConfigProvider
}

func NewHealthCheckHandler(logger *slog.Logger, cfg config.ConfigProvider) *HealthCheckHandler {
	return &HealthCheckHandler{
		Logger: logger,
		Config: cfg,
	}
}

func (h *HealthCheckHandler) Check(w http.ResponseWriter, r *http.Request) {
	mp := helpers.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.Config.GetEnv(),
		},
	}
	err := helpers.WriteJSON(w, http.StatusOK, mp, nil)
	if err != nil {
		// TODO: handle error
		return
	}
}
