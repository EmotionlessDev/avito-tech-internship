package models

type TeamMember struct {
	UserID   string `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	IsActive bool   `db:"is_active" json:"is_active"`
	TeamName string `db:"team_name" json:"-"`
}
