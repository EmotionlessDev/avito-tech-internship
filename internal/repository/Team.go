package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
)

var ErrTeamNotFound error = errors.New("team_name not found")
var ErrTeamDuplicate error = errors.New("team_name already exists")

type TeamRepository interface {
	CreateTeam(ctx context.Context, name string) error
	GetByName(ctx context.Context, name string) (*models.Team, error)
	AddMembers(ctx context.Context, teamName string, members []models.User) error
}

type TeamRepo struct {
	db      *sql.DB
	members TeamMemberRepository
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{
		db:      db,
		members: NewTeamMemberRepo(db),
	}
}

func (r *TeamRepo) CreateTeam(ctx context.Context, name string) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)", name).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrTeamDuplicate
	}
	_, err = r.db.ExecContext(ctx, "INSERT INTO teams (team_name) VALUES ($1)", name)
	return err
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (*models.Team, error) {
	var team models.Team
	team.TeamName = name

	members, err := r.members.ListByTeam(ctx, name)
	if err != nil {
		return nil, err
	}
	team.Members = members

	if len(members) == 0 {
		var exists bool
		err = r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)", name).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrTeamNotFound
		}
	}

	return &team, nil
}

func (r *TeamRepo) AddMembers(ctx context.Context, teamName string, members []models.User) error {
	return r.members.AddMembers(ctx, teamName, members)
}
