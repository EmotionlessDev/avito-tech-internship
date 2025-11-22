package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
	"github.com/EmotionlessDev/avito-tech-internship/internal/services"
)

var errMissingTeamName = errors.New("missing team_name parameter")

type TeamHandler struct {
	ErrorResponder *ErrorResponder
	Logger         *slog.Logger
	TeamService    *services.TeamService
}

func NewTeamHandler(
	ErrorResponder *ErrorResponder,
	logger *slog.Logger,
	teamService *services.TeamService,
) *TeamHandler {
	return &TeamHandler{
		ErrorResponder: ErrorResponder,
		Logger:         logger,
		TeamService:    teamService,
	}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TeamName    string        `json:"team_name"`
		TeamMembers []models.User `json:"members"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		h.ErrorResponder.BadRequest(w, r, err)
		return
	}

	team := models.Team{
		TeamName: input.TeamName,
		Members:  input.TeamMembers,
	}

	err = h.TeamService.CreateTeamWithMembers(r.Context(), team)
	if err != nil {
		if errors.Is(err, repository.ErrTeamDuplicate) {
			h.ErrorResponder.Conflict(w, r, "team already exists")
			return
		}

		h.ErrorResponder.InternalServerError(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusCreated, helpers.Envelope{"team": team}, nil)
	if err != nil {
		h.ErrorResponder.ServerError(w, r, err)
		return
	}
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	if teamName == "" {
		h.ErrorResponder.BadRequest(w, r, errMissingTeamName)
		return
	}

	team, err := h.TeamService.GetTeam(r.Context(), teamName)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			h.ErrorResponder.NotFound(w, r)
			return
		}
		h.ErrorResponder.InternalServerError(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"team": team}, nil)
	if err != nil {
		h.ErrorResponder.ServerError(w, r, err)
		return
	}
}
