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
	Create(ctx context.Context, name string) error
	GetByName(ctx context.Context, name string) (*models.Team, error)
}

type TeamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, teamName string) error {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)`,
		teamName,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrTeamDuplicate
	}

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO teams (team_name) VALUES ($1)`,
		teamName,
	)

	return err
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (*models.Team, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)`,
		name,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTeamNotFound
	}

	team := &models.Team{
		TeamName: name,
		Members:  []models.User{},
	}

	rows, err := r.db.QueryContext(ctx, `
        SELECT user_id, username, is_active
        FROM users
        WHERE team_name = $1
    `, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.IsActive); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, u)
	}

	return team, nil
}
