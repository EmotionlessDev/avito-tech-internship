package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/EmotionlessDev/avito-tech-internship/internal/helpers"
	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
)

var errMissingTeamName = errors.New("missing team_name parameter")

type TeamHandler struct {
	TeamRepo       repository.TeamRepository
	TeamMemberRepo repository.TeamMemberRepository
	ErrorResponder *ErrorResponder
	Logger         *slog.Logger
}

func NewTeamHandler(
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	ErrorResponder *ErrorResponder,
	logger *slog.Logger,
) *TeamHandler {
	return &TeamHandler{
		TeamRepo:       teamRepo,
		TeamMemberRepo: teamMemberRepo,
		ErrorResponder: ErrorResponder,
		Logger:         logger,
	}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TeamName string `json:"team_name"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		h.ErrorResponder.BadRequest(w, r, err)
		return
	}

	team := &models.Team{
		TeamName: input.TeamName,
	}

	err = h.TeamRepo.Create(r.Context(), team)
	if err != nil {
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

	team, err := h.TeamRepo.GetByName(r.Context(), teamName)
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
