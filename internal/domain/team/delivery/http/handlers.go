package http

import (
	"context"
	"errors"
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
		common.MethodNotAllowedResponse(w)
		return
	}

	var req AddTeamRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		common.BadRequestResponse(w, err)
		return
	}

	ctx := context.Background()
	teamEntity := &team.Team{Name: req.TeamName}

	domainMembers := make([]team.User, 0, len(req.Members))
	members := make([]Member, 0, len(req.Members))
	for _, m := range req.Members {
		domainMembers = append(domainMembers, team.User{
			ID:       m.ID,
			Name:     m.Name,
			TeamName: req.TeamName,
			IsActive: m.IsActive,
		})
		members = append(members, Member{
			ID:       m.ID,
			Name:     m.Name,
			IsActive: m.IsActive,
		})
	}

	err := h.addService.Add(ctx, teamEntity, domainMembers)
	if err != nil {

		if errors.Is(err, common.ErrTeamDuplicate) {
			common.BadRequestResponse(w, common.ErrTeamDuplicate)
			return
		}

		common.InternalServerErrorResponse(w, err)
		return
	}

	rsp := &AddTeamResponse{
		Team:    teamEntity.Name,
		Members: members,
	}

	helpers.WriteJSONObj(w, http.StatusCreated, rsp, nil)
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.MethodNotAllowedResponse(w)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		common.BadRequestResponse(w, common.ErrInvalidTeamName)
		return
	}

	ctx := r.Context()

	res, err := h.getService.Get(ctx, teamName)
	if err != nil {
		if errors.Is(err, common.ErrTeamNotFound) {
			common.NotFoundResponse(w, common.ErrTeamNotFound)
			return
		}

		common.InternalServerErrorResponse(w, err)
		return
	}

	rsp := &GetTeamResponse{
		Team:    res.Team.Name,
		Members: res.Members,
	}

	helpers.WriteJSONObj(w, http.StatusOK, rsp, nil)
}
