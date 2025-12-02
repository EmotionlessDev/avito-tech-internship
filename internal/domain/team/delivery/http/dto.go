package http

import (
	"strconv"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
	"github.com/EmotionlessDev/avito-tech-internship/internal/validator"
)

type AddTeamResponse struct {
	Team    string   `json:"team_name"`
	Members []Member `json:"members"`
}

type GetTeamResponse struct {
	Team    string      `json:"team_name"`
	Members []team.User `json:"members"`
}

type AddTeamMemberRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type AddTeamRequest struct {
	TeamName string                 `json:"team_name"`
	Members  []AddTeamMemberRequest `json:"members"`
}

type Member struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (t AddTeamRequest) Validate(v *validator.Validator) {
	v.Check(t.TeamName != "", "team_name", "team_name is required")
	for i, m := range t.Members {
		if m.ID == "" {
			v.Check(false, "members["+strconv.Itoa(i)+"].id", "member id is required")
		}
		if m.Name == "" {
			v.Check(false, "members["+strconv.Itoa(i)+"].name", "member name is required")
		}
	}
}
