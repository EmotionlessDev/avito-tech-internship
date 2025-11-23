package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/service/add"
	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
)

type Handler struct {
	Service *add.Service
}

func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.WriteJSON(w, http.StatusMethodNotAllowed, helpers.Envelope{
			"error": "method not allowed",
		}, nil)
		return
	}

	var req AddTeamRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error":   "invalid_request",
			"message": err.Error(),
		}, nil)
		return
	}

	ctx := context.Background()
	teamEntity := &team.Team{Name: req.TeamName}

	err := h.Service.Add(ctx, teamEntity, req.Members)
	if err != nil {

		if errors.Is(err, common.ErrTeamDuplicate) {
			helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
				"error":   "team_exists",
				"message": fmt.Sprintf("%s already exists", req.TeamName),
			}, nil)
			return
		}

		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{
			"error":   "internal_error",
			"message": err.Error(),
		}, nil)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, helpers.Envelope{
		"team":    teamEntity.Name,
		"members": req.Members,
	}, nil)
}
