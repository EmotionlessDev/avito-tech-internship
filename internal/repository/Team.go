package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`, team.TeamName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check team existence: %w", err)
	}
	if exists {
		return ErrTeamDuplicate
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO teams (team_name) VALUES ($1)`, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to insert team: %w", err)
	}

	for _, m := range team.TeamMembers {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO team_members (user_id, username, team_name, is_active)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) DO UPDATE 
                SET username = EXCLUDED.username,
                    team_name = EXCLUDED.team_name,
                    is_active = EXCLUDED.is_active
        `, m.UserID, m.Username, team.TeamName, m.IsActive)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (*models.Team, error) {
	var team models.Team
	team.TeamMembers = []models.TeamMember{}

	rows, err := r.db.QueryContext(ctx, `
    SELECT t.team_name, m.user_id, m.username, m.is_active
    FROM teams t
    LEFT JOIN team_members m ON t.team_name = m.team_name
    WHERE t.team_name = $1
`, name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m models.TeamMember
		if err := rows.Scan(&team.TeamName, &m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		m.TeamName = team.TeamName
		team.TeamMembers = append(team.TeamMembers, m)
	}

	return &team, nil
}
