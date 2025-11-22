package models

type Team struct {
	TeamName    string       `db:"team_name" json:"team_name"`
	TeamMembers []TeamMember `db:"team_members" json:"team_members"`
}
