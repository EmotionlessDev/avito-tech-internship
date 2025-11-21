package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
)

var ErrTeamNotFound error = errors.New("team not found")

type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	GetByName(ctx context.Context, name string) (*models.Team, error)
}

type TeamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, team *models.Team) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO teams (team_name) VALUES ($1)`,
		team.TeamName,
	)
	return err
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (*models.Team, error) {
	row := r.db.QueryRowContext(ctx, `SELECT team_name FROM teams WHERE team_name = $1`, name)

	var team models.Team
	err := row.Scan(&team.TeamName)
	if err == sql.ErrNoRows {
		return nil, ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team by name: %w", err)
	}

	return &team, nil
}
