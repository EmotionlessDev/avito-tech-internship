package http

import (
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
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
