package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/service/add"
)

type Handler struct {
	Service *add.Service
}

func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	t := &team.Team{Name: req.TeamName}
	ctx := context.Background()

	err := h.Service.Add(ctx, t, req.Members)
	if err != nil {
		if errors.Is(err, common.ErrTeamDuplicate) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Code:    "TEAM_EXISTS",
				Message: fmt.Sprintf("%s already exists", req.TeamName),
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AddTeamResponse{
		Team:    t.Name,
		Members: req.Members,
	})
}
