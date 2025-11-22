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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var exists bool
	err = tx.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)`,
		team.TeamName,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrTeamDuplicate
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO teams (team_name) VALUES ($1)`,
		team.TeamName,
	)
	if err != nil {
		return err
	}

	for _, u := range team.Members {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET 
				username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active
		`, u.UserID, u.Username, team.TeamName, u.IsActive)

		if err != nil {
			return err
		}
	}

	return nil
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

	team := models.Team{TeamName: name}

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

	return &team, nil
}
