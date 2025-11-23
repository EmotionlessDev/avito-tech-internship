package http

import "github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"

type AddTeamRequest struct {
	TeamName string      `json:"team_name"`
	Members  []team.User `json:"members"`
}

type AddTeamResponse struct {
	Team    string      `json:"team_name"`
	Members []team.User `json:"members"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
