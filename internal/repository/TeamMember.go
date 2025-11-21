package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
)

var ErrTeamMemberNotFound error = errors.New("team member not found")

type TeamMemberRepository interface {
	Create(ctx context.Context, member *models.TeamMember) error
	ListByTeam(ctx context.Context, teamName string) ([]models.TeamMember, error)
}

type TeamMemberRepo struct {
	db *sql.DB
}

func NewTeamMemberRepo(db *sql.DB) *TeamMemberRepo {
	return &TeamMemberRepo{db: db}
}

func (r *TeamMemberRepo) Create(ctx context.Context, m *models.TeamMember) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (user_id, username, is_active, team_name)
         VALUES ($1, $2, $3, $4)`,
		m.UserID, m.Username, m.IsActive, m.TeamName,
	)
	return err
}

func (r *TeamMemberRepo) ListByTeam(ctx context.Context, teamName string) ([]models.TeamMember, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT user_id, username, is_active, team_name FROM users WHERE team_name = $1`,
		teamName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query team members: %w", err)
	}
	defer func() {
		_ = rows.Close() // ignore error on close
	}()

	var members []models.TeamMember
	for rows.Next() {
		var m models.TeamMember
		err = rows.Scan(&m.UserID, &m.Username, &m.IsActive, &m.TeamName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team member: %w", err)
		}
		members = append(members, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return members, nil
}
