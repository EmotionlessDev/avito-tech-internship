package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team/service/get"
	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
)

type AddService interface {
	Add(ctx context.Context, t *team.Team, members []team.User) error
}

type GetService interface {
	Get(ctx context.Context, teamName string) (*get.TeamWithMembers, error)
}

type Handler struct {
	addService AddService
	getService GetService
}

func NewHandler(addService AddService, getService GetService) *Handler {
	return &Handler{
		addService: addService,
		getService: getService,
	}
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

	members := make([]team.User, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, team.User{
			ID:       m.ID,
			Name:     m.Name,
			TeamName: req.TeamName,
			IsActive: m.IsActive,
		})
	}

	err := h.addService.Add(ctx, teamEntity, members)
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
		"team_name": teamEntity.Name,
		"members":   req.Members,
	}, nil)
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.WriteJSON(w, http.StatusMethodNotAllowed, helpers.Envelope{
			"error": "method not allowed",
		}, nil)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{
			"error": "team_name is required",
		}, nil)
		return
	}

	ctx := r.Context()

	res, err := h.getService.Get(ctx, teamName)
	if err != nil {
		if errors.Is(err, common.ErrTeamNotFound) {
			helpers.WriteJSON(w, http.StatusNotFound, helpers.Envelope{
				"error": "team not found",
			}, nil)
			return
		}

		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{
			"error": err.Error(),
		}, nil)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{
		"team_name": res.Team.Name,
		"members":   res.Members,
	}, nil)
}
