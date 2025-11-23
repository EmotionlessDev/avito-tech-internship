package http

type SetUserActiveRequest struct {
	ID       string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetUserActiveResponse struct {
	ID       string `json:"user_id"`
	IsActive bool   `json:"is_active"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
}
